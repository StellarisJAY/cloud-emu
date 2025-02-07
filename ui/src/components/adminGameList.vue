<template>
  <a-card :bordered="false" class="center-card">
    <template #extra>
      <a-button type="primary" @click="showUploadModal">上传</a-button>
    </template>
    <a-table :columns="columns" size="small"
             :data-source="games" :pagination="false" style="width: 100%"></a-table>
    <a-pagination v-model:current="query.page" :total="total" v-model:page-size="query.pageSize" @change="onPaginationChange"></a-pagination>
  </a-card>
  <a-modal :open="uploadModalOpen" @cancel="onUploadClose" @ok="uploadFile">
    <a-form>
      <a-form-item label="游戏名">
        <a-input v-model:value="gameName"></a-input>
      </a-form-item>
      <a-form-item label="模拟器">
        <a-select :options="emulatorTypeOptions" v-model:value="selectedEmulatorType"></a-select>
      </a-form-item>
      <a-upload-dragger :multiple="false"
                        :file-list="fileList"
                        @change="onUploadChange"
                        @drop="onUploadDrop" :before-upload="() => {return false}">
        <p class="ant-upload-text">点击或拖拽文件到该区域来上传</p>
      </a-upload-dragger>
    </a-form>
  </a-modal>
</template>

<script>
import {Row, Col, List, Card, Modal, UploadDragger, message} from "ant-design-vue";
import { Button, Table, Form, Select, Input, Pagination} from "ant-design-vue";
import api from "../api/request.js";
import emulatorAPI from "../api/emulator.js";

export default {
  components: {
    ARow: Row,
    ACol: Col,
    AButton: Button,
    AList: List,
    AListItem: List.Item,
    ACard: Card,
    ATable: Table,
    AModal: Modal,
    AUploadDragger: UploadDragger,
    AFormItem: Form.Item,
    ASelect: Select,
    AForm: Form,
    AInput: Input,
    APagination: Pagination,
  },
  data() {
    return {
      games: [],
      query: {
        emulatorId: 0,
        gameName: "",
        page: 1,
        pageSize: 10,
      },
      total: 0,
      columns: [
        {title: "游戏名称", dataIndex: "gameName"},
        {title: "模拟器类型", dataIndex: "emulatorType"},
        {title: "大小", dataIndex: "size"},
      ],

      uploadModalOpen: false,
      fileList: [],
      gameName: "",
      emulatorTypes: [],
      emulatorTypeOptions: [],
      selectedEmulatorType: null,
    }
  },
  created() {
    this.listGames();
    emulatorAPI.listEmulatorTypes().then(resp=>{
      this.emulatorTypes = resp.data;
      this.selectedEmulatorType = this.emulatorTypes[0];
      let options = [];
      this.emulatorTypes.forEach(emulatorType => {
        options.push({label: emulatorType, value: emulatorType});
      })
      this.emulatorTypeOptions = options;
    });
  },
  unmounted() {
  },
  methods: {
    onPaginationChange(page) {
      this.query.page = page;
      this.listGames();
    },
    listGames() {
      api.get("/game", this.query).then(res => {
        this.games = res.data;
        this.total = res["total"];
        this.games.forEach(game => {
          game.size = this.formatGameSize(game.size);
        });
      });
    },
    formatGameSize(size) {
      if (size <= 1000) {
        return size + "B"
      }else if (size <= 1000000) {
        return Math.ceil(size / 1000) + "KB";
      }else {
        return Math.ceil(size/10000000) + "MB";
      }
    },
    showUploadModal() {
      this.uploadModalOpen = true;
    },
    onUploadChange(e) {
      const file = e.file;
      if (file.status !== "removed") {
        this.fileList = [e.file];
      } else {
        this.fileList = [];
      }
    },
    onUploadDrop(e) {
    },
    onUploadClose() {
      this.fileList = [];
      this.uploadModalOpen = false;
    },
    uploadFile() {
      const formData = new FormData();
      formData.append("gameName", this.gameName);
      formData.append("emulatorType", this.selectedEmulatorType);
      formData.append("file", this.fileList[0]);
      api.postForm("/game/upload", formData).then(res => {
        message.success("上传成功");
        this.listGames();
      }).catch(resp=>{message.error(resp.message)});
    },
  }
}
</script>

<style scoped>
</style>