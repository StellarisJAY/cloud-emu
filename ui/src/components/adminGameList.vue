<template>
  <a-card :bordered="false" class="center-card">
    <template slot="header">
      <a-button>上传</a-button>
    </template>
    <a-table :columns="columns"
             :data-source="games"
             :pagination="{pageSize: query.pageSize, current: query.page}"></a-table>
  </a-card>
</template>

<script>
import {Row, Col, List, Card} from "ant-design-vue";
import { Button, Table} from "ant-design-vue";
import api from "../api/request.js";

export default {
  components: {
    ARow: Row,
    ACol: Col,
    AButton: Button,
    AList: List,
    AListItem: List.Item,
    ACard: Card,
    ATable: Table,
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
      columns: [
        {title: "游戏名称", dataIndex: "gameName"},
        {title: "模拟器类型", dataIndex: "emulatorType"},
        {title: "大小", dataIndex: "size"},
        {title: "上传时间", dataIndex: "addTime"},
      ],
    }
  },
  created() {
    this.listGames();
  },
  unmounted() {
  },
  methods: {
    listGames() {
      api.get("/game", this.query).then(res => {
        this.games = res.data;
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
    }
  }
}
</script>

<style scoped>
</style>