<script setup>
</script>

<template>
  <a-card :bordered="false" class="center-card">
    <a-form>
      <a-form-item label="用户名">
        <a-input v-model:value="userDetail.userName" disabled></a-input>
      </a-form-item>
      <a-form-item label="昵称">
        <a-input v-model:value="userDetail.nickName"></a-input>
      </a-form-item>
      <a-form-item label="密码">
        <a-input-password v-model:value="userDetail.password"></a-input-password>
      </a-form-item>
      <a-form-item label="邮箱">
        <a-input v-model:value="userDetail.email" disabled></a-input>
      </a-form-item>
      <a-button type="primary">修改</a-button>
    </a-form>
    <a-form v-if="userDetail.status === 3">
      <a-form-item label="激活码">
        <a-input v-model:value="activationCode"></a-input>
      </a-form-item>
      <a-button type="primary" @click="activateAccount">激活账号</a-button>
    </a-form>
  </a-card>
</template>

<script>
import {Row, Col, Card, message} from "ant-design-vue";
import { Button, Form, Input } from "ant-design-vue";
import userAPI from "../api/user.js";

export default {
  components: {
    ARow: Row,
    ACol: Col,
    AButton: Button,
    ACard: Card,
    AForm: Form,
    AFormItem: Form.Item,
    AInput: Input,
    AInputPassword: Input.Password,
  },
  data() {
    return {
      userDetail: {
        userId: 0,
        userName: "",
        addTime: "",
        nickName: "",
        status: 1,
        email: "",
        phone: "",
      },
      activationCode: "",
    }
  },
  created() {
    this.getUserDetail();
  },
  unmounted() {
  },
  methods: {
    getUserDetail() {
      userAPI.getLoginUserDetail().then(resp=>{
        this.userDetail = resp.data;
      });
    },
    activateAccount() {
      userAPI.activateAccount(this.activationCode)
          .then(resp=>{message.success("激活成功");this.getUserDetail();})
          .catch(_=>message.error("激活失败"));
    },
  }
}
</script>

<style scoped>
</style>