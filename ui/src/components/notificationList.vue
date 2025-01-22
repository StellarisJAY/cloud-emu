<script setup>
const NoticeTypeInvite = 2;
</script>

<template>
  <a-card :bordered="false" class="center-card">
    <template #extra>
      <a-button danger @click="clearInbox">清空</a-button>
    </template>
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
        <template #actions>
          <a-button type="link" @click="deleteNotification(item['notificationId'])">删除</a-button>
        </template>
      </a-list-item>
    </a-list>
    <a-pagination v-model:current="page" :total="total" :page-size="pageSize" @change="listNotifications"></a-pagination>
  </a-card>
</template>

<script>
import {Row, Col, List, Card, message, Pagination} from "ant-design-vue";
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
    APagination: Pagination,
  },
  data() {
    return {
      notifications: [],
      page: 1,
      pageSize: 10,
      total: 0,
    }
  },
  created() {
    this.listNotifications();
  },
  unmounted() {
  },
  methods: {
    listNotifications() {
      api.get("/inbox/notifications", {page: this.page, pageSize: this.pageSize}).then(resp=>{
        this.notifications = resp.data;
        this.total = resp["total"];
      });
    },
    getRoomLink(content) {
      const room = JSON.parse(content);
      return "/room/" + room["roomId"];
    },
    getRoomName(content) {
      const room = JSON.parse(content);
      return room["roomName"];
    },
    deleteNotification(id) {
      api.post("/inbox/notifications", {notificationIds: [id]}).then(resp=>{
        message.success(resp["message"]);
        this.listNotifications();
      }).catch(resp=>{
        message.error(resp["message"]);
      });
    },
    clearInbox() {
      api.delete("/inbox/clear").then(resp=>{
        message.success(resp["message"]);
        this.listNotifications();
      }).catch(resp=>{
        message.error(resp["message"]);
      });
    }
  }
}
</script>

<style scoped>
</style>