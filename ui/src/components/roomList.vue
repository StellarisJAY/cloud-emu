<template>
  <a-card :bordered="false" class="center-card">
    <template #extra>
      <a-button v-if="joined" type="primary" @click="openCreateRoomModal">新建房间</a-button>
    </template>
    <a-row>
      <a-form layout="inline">
        <a-col :span="7">
          <a-form-item label="房间名称">
            <a-input v-model:value="listRoomQuery.roomName"></a-input>
          </a-form-item>
        </a-col>
        <a-col :span="7">
          <a-form-item label="房主名称">
            <a-input v-model:value="listRoomQuery.hostName"></a-input>
          </a-form-item>
        </a-col>
        <a-col :span="7">
          <a-form-item label="访问权限">
            <a-select :options="joinTypeQueryOptions" v-model:value="listRoomQuery.joinType"></a-select>
          </a-form-item>
        </a-col>
        <a-col :span="3">
          <a-button @click="listRoom">查询</a-button>
        </a-col>
      </a-form>
    </a-row>
    <a-list :grid="{ gutter: 16, xs: 1, sm: 2, md: 2, lg: 2, xl: 3, xxl: 3 }" :data-source="rooms">
      <template #renderItem="{ item }">
        <a-list-item>
          <a-card :title="item.roomName">
            <template #extra v-if="joined">
              <a-button v-if="item['isHost']" danger @click="confirmDelete(item)">删除</a-button>
              <a-button v-else danger>退出</a-button>
            </template>
            <template #actions>
              <RouterLink v-if="joined" :to="'/room/' + item.roomId">进入</RouterLink>
              <a-button v-else type="link" @click="joinRoom(item)">加入</a-button>
            </template>
            <ul style="text-align: left">
              <li>模拟器：{{item["emulatorId"]==='0'?'无':item["emulatorName"]}}</li>
              <li>游戏：{{item["gameId"]==='0'?'无':item["gameName"]}}</li>
              <li>房主：{{ item["hostName"] }}</li>
              <li>人数：{{ item["memberCount"] }}/{{ item["memberLimit"] }}</li>
              <li>访问权限：{{ getJoinTypeName(item["joinType"]) }}</li>
              <li>创建时间：{{ item["addTime"] }}</li>
            </ul>
          </a-card>
        </a-list-item>
      </template>
    </a-list>

    <a-pagination v-model:current="listRoomQuery.page" v-model:pageSize="listRoomQuery.pageSize" :total="total" @change="onPageChange">
    </a-pagination>

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
    <a-modal v-if="joined" :open="deleteRoomModalOpen" title="警告" @cancel="closeConfirmDelete" @ok="deleteRoom">
      <p>删除房间会同时永久删除所有存档文件，你可以先将存档转移到其他房间</p>
      <p>请在下方输入框输入房间名：{{deletedRoom["roomName"]}}</p>
      <a-input v-model:value="deleteRoomConfirmInput"></a-input>
    </a-modal>
  </a-card>
</template>

<script>
import {
  Card,
  Button,
  List,
  Modal,
  Form,
  Input,
  Switch,
  Select,
  Descriptions,
  InputNumber,
  Pagination
} from 'ant-design-vue';
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
    APagination: Pagination,
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
      selectedJoinType: 1,
      joinTypeOptions: [],

      joinTypeQueryOptions: [],
      listRoomQuery: {
        roomName: "",
        hostName: "",
        joinType: 0,
        page: 1,
        pageSize: 10
      },
      total: 0,
      deleteRoomModalOpen: false,
      deletedRoom: null,
      deleteRoomConfirmInput: "",
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
    onPageChange(page) {
      this.listRoomQuery.page = page;
      this.listRoom();
    },
    listRoom(e) {
      if (this.joined) {
        this.listJoinedRooms()
      } else {
        this.listAllRooms()
      }
    },
    listJoinedRooms() {
      roomAPI.listJoinedRooms(this.listRoomQuery).then(resp=>{
        this.rooms = resp.data;
        this.total = resp.total;
      });
    },
    listAllRooms() {
      roomAPI.listAllRooms(this.listRoomQuery).then(resp=>{
        this.rooms = resp.data;
        this.total = resp.total;
      });
    },
    createRoom() {
      if (this.createRoomState.name === "") {
        message.warn("请输入房间名");
        return;
      }
      if (this.createRoomState.joinType === 2) {
        if (this.createRoomState.password === "") {
          message.warn("请输入房间密码");
          return;
        }
      }
      const _this = this;
      api.post("/room", this.createRoomState).then(_ => {
        _this.listJoinedRooms();
        _this.createRoomModalOpen = false
        message.success("创建成功");
      });
    },
    joinRoom(room) {
      router.push(`/room/`+room["roomId"]);
    },
    leaveRoom(id) {
      // TODO leave room
    },
    openCreateRoomModal: function() {
      this.joinTypeOptions = configs.getEnumOptions("roomJoinTypeEnum");
      this.createRoomModalOpen = true;
    },
    getJoinTypeName(id) {
      return configs.getEnumName("roomJoinTypeEnum", id);
    },
    confirmDelete(room) {
      this.deletedRoom = room;
      this.deleteRoomModalOpen = true;
    },
    closeConfirmDelete() {
      this.deleteRoomModalOpen = false;
    },
    deleteRoom() {
      if (this.deleteRoomConfirmInput !== this.deletedRoom["roomName"]) {
        message.warn("请正确输入房间名称");
        return;
      }
      api.delete("/room/"+this.deletedRoom["roomId"]).then(_=>{
        message.success("删除成功");
        this.listRoom();
        this.closeConfirmDelete();
      }).catch(resp=>{
        message.error(resp["message"]);
      });
    },
  }
}
</script>