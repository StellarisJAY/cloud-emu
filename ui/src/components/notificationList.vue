<script setup>
const NoticeTypeInvite = 2;
</script>

<template>
  <a-card :bordered="false" class="center-card">
    <a-list>
      <a-list-item v-for="item in notifications">
        <a-list-item-meta v-if="item['type'] === NoticeTypeInvite">
          <template #title>
            {{item["senderNickName"]}}
          </template>
          <template #description>
            邀请您加入房间：<router-link :to="getRoomLink(item['content'])">{{getRoomName(item["content"])}}</router-link>
          </template>
        </a-list-item-meta>
      </a-list-item>
    </a-list>
  </a-card>
</template>

<script>
import {Row, Col, List, Card} from "ant-design-vue";
import { Button } from "ant-design-vue";
import api from "../api/request.js";

export default {
  components: {
    ARow: Row,
    ACol: Col,
    AButton: Button,
    AList: List,
    AListItem: List.Item,
    ACard: Card,
    AListItemMeta: List.Item.Meta,
  },
  data() {
    return {
      notifications: [],
      page: 1,
      pageSize: 10,
    }
  },
  created() {
    api.get("/inbox/notifications", {page: this.page, pageSize: this.pageSize}).then(resp=>{
      this.notifications = resp.data;
    });
  },
  unmounted() {
  },
  methods: {
    getRoomLink(content) {
      const room = JSON.parse(content);
      return "/room/" + room["roomId"];
    },
    getRoomName(content) {
      const room = JSON.parse(content);
      return room["roomName"];
    }
  }
}
</script>

<style scoped>
</style>