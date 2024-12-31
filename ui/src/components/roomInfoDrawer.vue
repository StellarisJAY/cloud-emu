<template>
  <div>
    <a-form layout="horizontal" labelAlign="right">
      <a-form-item label="房间名称">
        <a-input v-model:value="updateRoomForm.roomName" :disabled="!currentUserIsHost"></a-input>
      </a-form-item>
      <a-form-item label="访问权限">
        <a-select :options="joinTypeOptions" v-model:value="updateRoomForm.joinType" :disabled="!currentUserIsHost"></a-select>
      </a-form-item>
      <a-form-item label="密码" v-if="updateRoomForm.joinType === 2">
        <a-input-password v-model:value="updateRoomForm.password" :disabled="!currentUserIsHost"></a-input-password>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" @click="updateRoom" :disabled="updateRoomBtnDisabled || !currentUserIsHost">修改</a-button>
      </a-form-item>
    </a-form>
    <a-divider orientation="left">房间成员</a-divider>
    <a-list item-layout="vertical" :data-source="roomMembers">
      <template #renderItem="{ item }">
        <a-list-item>
          <a-row>
            <a-col :span="8">
              {{ item.nickName }}
              <a-tooltip>
                <template #title v-if="item.online">在线</template>
                <template #title v-else>离线</template>
                <CheckCircleTwoTone v-if="item.online"/>
                <CloseCircleTwoTone v-else />
              </a-tooltip>
              <a-tooltip title="房主" v-if="item.role===1">
                <CrownTwoTone/>
              </a-tooltip>
            </a-col>
            <a-col :span="12">
            </a-col>
            <a-col :span="4">
              <a-button :disabled="!currentUserIsHost || item.role === 1">删除</a-button>
            </a-col>
          </a-row>
        </a-list-item>
      </template>
    </a-list>
    <a-button @click="openUserPicker" type="primary" :disabled="!currentUserIsHost">邀请用户</a-button>
    <UserPicker :visible="userPickerOpen"/>
  </div>
</template>

<script>
import {
  Button,
  Col,
  Drawer,
  Form,
  Input,
  List,
  Row,
  Select,
  Divider, message, Tooltip,
} from 'ant-design-vue';
import {CrownTwoTone, CheckCircleTwoTone, CloseCircleTwoTone} from "@ant-design/icons-vue"
import constants from "../api/const.js";
import roomAPI from "../api/room.js";
import roomMemberAPI from "../api/roomMember.js";
import UserPicker from "./userPicker.vue";

export default {
  props: {
    fullRoomInfo: Object,
    rtcSession: Object,
    roomId: String,
  },
  components: {
    AButton: Button,
    ADrawer: Drawer,
    AForm: Form,
    AInput: Input,
    ARow: Row,
    ACol: Col,
    AList: List,
    AListItem: List.Item,
    CrownTwoTone,
    AInputPassword: Input.Password,
    AFormItem: Form.Item,
    ASelect: Select,
    ADivider: Divider,
    UserPicker: UserPicker,
    CheckCircleTwoTone,
    ATooltip: Tooltip,
    CloseCircleTwoTone,
  },
  data() {
    return {
      userRoomMember: null,
      roomMembers: [],
      joinTypeOptions: [],
      currentUserIsHost: true,
      updateRoomForm: {
        joinType: 1,
        password: null,
      },
      roomDetail: null,
      userPickerOpen: false,
      updateRoomBtnDisabled: false,
    }
  },
  created() {
    this.joinTypeOptions = constants.getEnumOptions("roomJoinTypeEnum");
    this.getUserRoomMember();
    this.listRoomMembers();
    this.getRoomDetail();
    addEventListener("memberDrawerOpen", _ => {
      this.getUserRoomMember();
      this.listRoomMembers();
      this.getRoomDetail();
    });
  },
  methods: {
    listRoomMembers: async function () {
      const resp = await roomMemberAPI.listRoomMember(this.roomId);
      this.roomMembers = resp.data;
    },

    updateRoom() {
      this.updateRoomBtnDisabled = true;
      roomAPI.updateRoom({
        "roomId": this.roomId,
        "roomName": this.updateRoomForm.roomName,
        "joinType": this.updateRoomForm.joinType,
        "password": this.updateRoomForm.password,
      }).then(() => {
        message.success("修改成功");
        this.updateRoomBtnDisabled = false;
      }).catch(resp=>{
        message.error(resp.message);
        this.updateRoomBtnDisabled = false;
      })
    },
    getRoomDetail() {
      roomAPI.getRoomDetail(this.roomId).then(resp=>{
        this.roomDetail = resp.data;
        this.updateRoomForm.joinType = this.roomDetail.joinType;
        this.updateRoomForm.roomName = this.roomDetail.roomName;
      });
    },
    getUserRoomMember() {
      roomMemberAPI.getUserRoomMember(this.roomId).then(resp => {
        this.userRoomMember = resp.data;
        this.currentUserIsHost = this.userRoomMember["role"] === 1;
      })
    },
    openUserPicker() {
      this.userPickerOpen = true;
      addEventListener("userPickerCancel", this.onUserPickerCancel);
      addEventListener("userPickerConfirm", this.onUserPickerConfirm);
    },
    inviteUsers(users) {
      users.forEach((user) => {
        roomMemberAPI.inviteRoomMember(this.roomId, user["userId"]);
      });
      message.success("已发送邀请");
    },
    onUserPickerCancel(ev) {
      this.userPickerOpen = false;
      removeEventListener("userPickerCancel", this.onUserPickerCancel);
    },
    onUserPickerConfirm(ev) {
      this.userPickerOpen = false;
      removeEventListener("userPickerConfirm", this.onUserPickerConfirm);
      const pickedUsers = ev.detail["users"];
      if (pickedUsers && pickedUsers.length > 0) {
        this.inviteUsers(pickedUsers);
      }
    },
  }
}
</script>