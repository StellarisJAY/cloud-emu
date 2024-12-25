<template>
  <a-card v-if="roomDetail['joinType'] === 2">
    <a-form>
      <a-form-item label="房间密码">
        <a-input-password v-model:value="password"></a-input-password>
      </a-form-item>
      <a-button type="primary" @click="joinRoom">加入</a-button>
    </a-form>
  </a-card>
</template>

<script>
import api from "../api/request.js";
import {Row, Col, message} from "ant-design-vue";
import { Button } from "ant-design-vue";
import { Form, FormItem, Modal, Input } from "ant-design-vue";
import roomAPI from "../api/room.js";
import router from "../router/index.js";

const RoomJoinTypePublic = 1;
const RoomJoinTypePassword = 2;
const RoomJoinTypeInvite = 3;

export default {
  components: {
    ARow: Row,
    ACol: Col,
    AButton: Button,
    AForm: Form,
    AFormItem: FormItem,
    AInput: Input,
    AModal: Modal,
    AInputPassword: Input.Password,
  },
  data() {
    return {
      roomId: null,
      roomDetail: null,
      password: null,
    }
  },
  created() {
    this.roomId = this.$route["params"]["roomId"];
    roomAPI.getRoomDetail(this.roomId).then(resp=>{
      this.roomDetail = resp["data"];
      if (this.roomDetail["joinType"] === RoomJoinTypeInvite || this.roomDetail["joinType"] === RoomJoinTypePublic) {
        this.joinRoom();
      }
    });
  },
  unmounted() {
  },
  methods: {
    joinRoom() {
      api.post("/room/join", {
        "roomId": this.roomId,
        "password": this.password,
      }).then(_=>{
        message.success("加入成功", 1).then(_=>router.back());
      }).catch(resp=>{
        if (resp["code"] === 403) {
          message.warn(resp["message"], 1).then(_=>router.go(-2));
        }else {
          message.error(resp["message"], 1).then(_=>router.go(-2));
        }
      });
    }
  }
}
</script>

<style scoped>
a-card {
  width: 50%;
  position: absolute;
  left: 25%;
}
</style>