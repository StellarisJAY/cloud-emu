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
    <a-pagination :page="page" :pageSize="pageSize" :total="totalPages" @change="onPaginationChange"></a-pagination>
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
      totalPages: 1,

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
        this.totalPages = Math.ceil(resp["total"] / this.pageSize);
      }).catch(_ => {
        message.error("获取存档列表失败");
      });
    },
    onPaginationChange: function (page) {
      this.page = page;
      this.listSaves();
    },
    loadSavedGame(id) {

    },
    deleteSavedGame(id) {

    },
  }
}
</script>