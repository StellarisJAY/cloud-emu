<template>
  <a-list item-layout="vertical" :data-source="saves">
    <template #renderItem="{ item }">
      <a-list-item>
        <a-descriptions :column="1">
          <a-descriptions-item label="模拟器">{{item["emulatorName"]}}</a-descriptions-item>
          <a-descriptions-item label="游戏">{{ item["gameName"] }}</a-descriptions-item>
          <a-descriptions-item label="时间">{{ item["addTime"]}}</a-descriptions-item>
        </a-descriptions>
        <a-button type="primary" @click="loadSavedGame(item['saveId'])" :disabled="loadSaveBtnDisable">加载</a-button>
        <a-button danger @click="deleteSavedGame(item['saveId'])" :disabled="loadSaveBtnDisable">删除</a-button>
      </a-list-item>
    </template>
    <a-pagination v-model:current="page" v-model:pageSize="pageSize" :total="total" @change="onPaginationChange"></a-pagination>
  </a-list>
</template>

<script>
import { List, Button, Descriptions, Pagination } from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import { message } from "ant-design-vue";
import api from "../api/request.js";

export default {
  props: {
    roomId: String,
  },
  components: {
    AButton: Button,
    ARow: Row,
    ACol: Col,
    AList: List,
    AListItem: List.Item,
    ADescriptions: Descriptions,
    ADescriptionsItem: Descriptions.Item,
    APagination: Pagination,
  },
  data() {
    return {
      saves: [],
      page: 1,
      pageSize: 10,
      total: 0,

      loadSaveBtnDisable: false,
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
  }
}
</script>