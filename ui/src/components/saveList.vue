<template>
  <a-list item-layout="vertical" :data-source="saves">
    <template #renderItem="{ item }">
      <a-list-item>
        <a-descriptions :column="1">
          <a-descriptions-item label="存档名">
            <a-input v-model:value="item['saveName']" @press-enter="renameSave(item)" :disabled="loadSaveBtnDisable"></a-input>
          </a-descriptions-item>
          <a-descriptions-item label="模拟器">{{item["emulatorName"]}}</a-descriptions-item>
          <a-descriptions-item label="游戏">{{ item["gameName"] }}</a-descriptions-item>
          <a-descriptions-item label="时间">{{ item["addTime"]}}</a-descriptions-item>
        </a-descriptions>
        <a-button type="primary" @click="loadSavedGame(item['saveId'])" :disabled="loadSaveBtnDisable">加载</a-button>
        <a-button type="primary" @click="openTransferSaveModal(item)" :disabled="loadSaveBtnDisable">转移</a-button>
        <a-button danger @click="deleteSavedGame(item['saveId'])" :disabled="loadSaveBtnDisable">删除</a-button>
      </a-list-item>
    </template>
    <a-pagination v-model:current="page" v-model:pageSize="pageSize" :total="total" @change="onPaginationChange"></a-pagination>
  </a-list>
  <a-modal :open="transferSaveModalOpen" title="转移存档" @cancel="_ => { transferSaveModalOpen = false }">
    <template #footer>
      <a-button @click="_ => { transferSaveModalOpen = false }">取消</a-button>
      <a-button type="primary" @click="transferSave" html-type="submit">确认</a-button>
    </template>
    <a-form-item label="目标房间">
      <a-select :options="roomOptions" v-model:value="selectedRoom"></a-select>
    </a-form-item>
  </a-modal>
</template>

<script>
import {
  List,
  Button,
  Descriptions,
  Pagination,
  Input,
  Modal,
  Form as AForm,
  InputNumber as AInputNumber, Select as ASelect
} from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import { message } from "ant-design-vue";
import api from "../api/request.js";
import roomAPI from "../api/room.js";

export default {
  props: {
    roomId: String,
  },
  components: {
    ASelect, AInputNumber, AForm,
    AButton: Button,
    ARow: Row,
    ACol: Col,
    AList: List,
    AListItem: List.Item,
    ADescriptions: Descriptions,
    ADescriptionsItem: Descriptions.Item,
    APagination: Pagination,
    AInput: Input,
    AModal: Modal,
    AFormItem: AForm.Item,
  },
  data() {
    return {
      saves: [],
      page: 1,
      pageSize: 10,
      total: 0,

      loadSaveBtnDisable: false,
      transferSaveModalOpen: false,

      rooms: [],
      roomOptions: [],
      selectedRoom: null,
      transferringSave: null,
    }
  },
  created() {
    this.listSaves();
    addEventListener("saveListOpen", _ => this.listSaves());
  },
  methods: {
    listSaves: function () {
      api.get("/game-save", {
        roomId: this.roomId,
        page: this.page,
        pageSize: this.pageSize,
      }).then(resp => {
        this.saves = resp.data;
        this.total = resp.total;
      }).catch(_ => {
        message.error("获取存档列表失败");
      });
    },
    onPaginationChange: function (page) {
      this.page = page;
      this.listSaves();
    },
    loadSavedGame(id) {
      this.loadSaveBtnDisable = true;
      api.post("/game-save/load", {"saveId": id,"roomId":this.roomId}).then(resp=>{
        this.loadSaveBtnDisable = false;
        message.success("读取存档成功");
      }).catch(resp=>{
        message.error(resp.message);
        this.loadSaveBtnDisable = false;
      });
    },
    deleteSavedGame(id) {
      this.loadSaveBtnDisable = true;
      api.delete("/game-save/" + id).then(resp=>{
        this.loadSaveBtnDisable = false;
        message.success(resp["message"]);
        this.listSaves();
      }).catch(resp=>{
        this.loadSaveBtnDisable = false;
        message.error(resp.message);
      })
    },
    renameSave(save) {
      api.put("/game-save/rename", {saveId: save.saveId, saveName: save["saveName"]}).then(resp=>{
        message.success(resp["message"]);
      }).catch(resp=>{
        message.error(resp.message);
      })
    },
    openTransferSaveModal(save) {
      roomAPI.listJoinedRooms({
        hostOnly: true,
        page: 1,
        pageSize: 20,
      }).then(resp=>{
        this.transferSaveModalOpen = true;
        this.rooms = resp.data;
        this.roomOptions = [];
        this.rooms.forEach(item=>{
          this.roomOptions.push({label:item["roomName"],value: item["roomId"]});
        });
        this.selectedRoom = this.rooms[0]["roomId"];
        this.transferringSave = save;
      }).catch(resp=>{
        message.error(resp.message);
      })
    },
    transferSave() {
      api.post("/game-save/transfer", {
        saveId: this.transferringSave["saveId"],
        roomId: this.selectedRoom,
      }).then(resp=>{
        message.success(resp["message"]);
      }).catch(resp=>{
        message.error(resp.message);
      });
    }
  }
}
</script>