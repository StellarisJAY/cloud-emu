<template>
  <div>
    <a-form layout="horizontal" labelAlign="right">
      <a-form-item label="访问权限">
        <a-select :options="joinTypeOptions" v-model:value="updateRoomForm.joinType" :disabled="!currentUserIsHost"></a-select>
      </a-form-item>
      <a-form-item label="密码" v-if="updateRoomForm.joinType === 2">
        <a-input-password v-model:value="updateRoomForm.password" :disabled="!currentUserIsHost"></a-input-password>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" @click="updateRoom" :disabled="!currentUserIsHost">确认</a-button>
      </a-form-item>
    </a-form>
    <a-divider orientation="left">房间成员</a-divider>
    <a-list item-layout="vertical" :data-source="roomMembers">
      <template #renderItem="{ item }">
        <a-list-item>
          <a-row>
            <a-col :span="8">
              <CrownTwoTone v-if="item.role === 1" />{{ item.nickName }}
            </a-col>
          </a-row>
        </a-list-item>
      </template>
    </a-list>
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
    Divider,
} from 'ant-design-vue';
import {CrownTwoTone} from "@ant-design/icons-vue"
import constants from "../api/const.js";
import roomAPI from "../api/room.js";
import roomMemberAPI from "../api/roomMember.js";

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

    },
    getRoomDetail() {
      roomAPI.getRoomDetail(this.roomId).then(resp=>{
        this.roomDetail = resp.data;
        this.updateRoomForm.joinType = this.roomDetail.joinType;
      });
    },
    getUserRoomMember() {
      roomMemberAPI.getUserRoomMember(this.roomId).then(resp => {
        this.userRoomMember = resp.data;
        this.currentUserIsHost = this.userRoomMember["role"] === 1;
      })
    }
  }
}
</script>