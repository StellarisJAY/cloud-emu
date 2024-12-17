<template>
  <a-card :bordered="false">
    <template #extra>
      <a-button v-if="joined" type="primary" @click="openCreateRoomModal">新建房间</a-button>
    </template>
    <a-row>
      <a-form layout="inline">
        <a-form-item label="房间名称">
          <a-input v-model:value="listRoomQuery.roomName"></a-input>
        </a-form-item>
        <a-form-item label="房主名称">
          <a-input v-model:value="listRoomQuery.hostName"></a-input>
        </a-form-item>
        <a-form-item label="访问权限">
          <a-select :options="joinTypeQueryOptions" v-model:value="listRoomQuery.joinType"></a-select>
        </a-form-item>
      </a-form>
      <a-button @click="listRoom">查询</a-button>
    </a-row>
    <a-list :grid="{ gutter: 16, xs: 1, sm: 2, md: 2, lg: 2, xl: 3, xxl: 3 }" :data-source="rooms">
      <template #renderItem="{ item }">
        <a-list-item>
          <a-card :title="item.roomName">
            <template #extra v-if="joined">
              <a-button v-if="item.role === 1" danger>删除</a-button>
              <a-button v-else danger>退出</a-button>
            </template>
            <template #actions>
              <RouterLink v-if="joined" :to="'/room/' + item.roomId">进入</RouterLink>
              <a-button v-else type="link" @click="joinRoom(item)">加入</a-button>
            </template>
            <ul style="text-align: left">
              <li>房主：{{ item["hostName"] }}</li>
              <li>人数：{{ item["memberCount"] }}/{{ item["memberLimit"] }}</li>
              <li>模拟器：{{ item["emulatorName"] }}</li>
              <li>访问权限：{{ getJoinTypeName(item["joinType"]) }}</li>
              <li>创建时间：{{ item["addTime"] }}</li>
            </ul>
          </a-card>
        </a-list-item>
      </template>
    </a-list>

    <a-modal v-if="joined" :open="createRoomModalOpen" title="新建房间" @cancel="_ => { createRoomModalOpen = false }">
      <template #footer>
        <a-button @click="_ => { createRoomModalOpen = false }">取消</a-button>
        <a-button type="primary" @click="createRoom()" html-type="submit">创建</a-button>
      </template>
      <a-form layout="vertical" :model="createRoomState" :label-col="{ span: 4 }">
        <a-form-item label="房间名" name="name" :rules="{ required: true, message: '请输入房间名' }">
          <a-input v-model:value="createRoomState.name"></a-input>
        </a-form-item>
        <a-form-item label="访问权限" name="joinType">
          <a-select :options="joinTypeOptions" v-model:value="createRoomState.joinType"></a-select>
        </a-form-item>
        <a-form-item label="密码" name="password" v-if="createRoomState.joinType === 2">
          <a-input-password v-model:value="createRoomState.password"></a-input-password>
        </a-form-item>
        <a-form-item label="人数" name="memberLimit">
          <a-input-number v-model:value="createRoomState.memberLimit" :max="4" :min="1"></a-input-number>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script>
import { Card, Button, List, Modal, Form, Input, Switch, Select, Descriptions, InputNumber} from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import { message } from "ant-design-vue";
import api from "../api/request.js";
import router from "../router/index.js";
import { RouterLink } from "vue-router";
import configs from "../api/const.js";
import roomAPI from "../api/room.js";

export default {
  props: {
    joined: Boolean
  },
  components: {
    ARow: Row,
    ACol: Col,
    ACard: Card,
    AButton: Button,
    AList: List,
    AListItem: List.Item,
    AModal: Modal,
    AForm: Form,
    AFormItem: Form.Item,
    AInput: Input,
    ASwitch: Switch,
    AInputPassword: Input.Password,
    ASelect: Select,
    ADescriptions: Descriptions,
    ADescriptionsItem: Descriptions.Item,
    AInputNumber: InputNumber,
  },
  data() {
    return {
      rooms: [
      ],
      createRoomModalOpen: false,
      createRoomState: {
        name: "",
        joinType: 1,
        password: "",
        memberLimit: 1
      },
      searchInput: "",
      joinRoomModalOpen: false,
      joinRoomFormState: {
        id: 0,
        password: ""
      },
      supportedEmulators: [],
      selectedEmulatorName: "",
      selectedEmulator: {},
      emulatorOptions: [],
      selectedJoinType: 1,
      joinTypeOptions: [],

      joinTypeQueryOptions: [],
      listRoomQuery: {
        roomName: "",
        hostName: "",
        joinType: 0,
        page: 1,
        pageSize: 10
      }
    }
  },
  created() {
    if (this.joined) {
      this.listJoinedRooms()
    } else {
      this.listAllRooms()
    }
    this.joinTypeQueryOptions = configs.getEnumOptionsWithAll("roomJoinTypeEnum");
  },
  methods: {
    listRoom() {
      if (this.joined) {
        this.listJoinedRooms()
      } else {
        this.listAllRooms()
      }
    },
    listJoinedRooms() {
      roomAPI.listJoinedRooms(this.listRoomQuery).then(resp=>this.rooms = resp.data);
    },
    listAllRooms() {
      roomAPI.listAllRooms(this.listRoomQuery).then(resp=>this.rooms = resp.data);
    },
    createRoom() {
      if (this.createRoomState.name === "") {
        message.warn("请输入房间名");
        return;
      }
      const _this = this;
      api.post("/room", this.createRoomState).then(_ => {
        _this.listJoinedRooms();
        _this.createRoomModalOpen = false
        message.success("创建成功");
      });
    },
    joinRoom(room) {

    },
    leaveRoom(id) {
      // TODO leave room
    },
    openCreateRoomModal: function() {
      this.joinTypeOptions = configs.getEnumOptions("roomJoinTypeEnum");
      this.createRoomModalOpen = true;
    },
    onSelectEmulatorChange: function() {
      this.selectedEmulator = this.supportedEmulators.find(emulator => emulator.name === this.selectedEmulatorName);
    },
    getJoinTypeName(id) {
      return configs.getEnumName("roomJoinTypeEnum", id);
    }
  }
}
</script>