<template>
  <div>
    <a-form layout="horizontal" labelAlign="right">
      <a-form-item label="模拟器">
        <a-select :options="emulatorOptions" v-model:value="updateForm.emulatorId"
                  @change="onEmulatorSelectChange" :disabled="!currentUserIsHost"></a-select>
      </a-form-item>
      <a-descriptions :column="1" v-if="selectedEmulator" title="模拟器详情">
        <a-descriptions-item label="类型">{{selectedEmulator["emulatorType"]}}</a-descriptions-item>
        <a-descriptions-item label="描述">{{ selectedEmulator["description"] }}</a-descriptions-item>
        <a-descriptions-item label="提供者">{{ selectedEmulator["provider"] }}</a-descriptions-item>
        <a-descriptions-item label="支持存档">{{ selectedEmulator["supportSave"] ? "是" : "否" }}</a-descriptions-item>
        <a-descriptions-item label="支持画面设置">{{ selectedEmulator["supportGraphicSetting"] ? "是" : "否"
          }}</a-descriptions-item>
      </a-descriptions>
      <a-form-item label="游戏">
        <a-select :options="emulatorGameOptions" v-model:value="updateForm.gameId" :disabled="!currentUserIsHost"></a-select>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" @click="restart" :disabled="!currentUserIsHost ">重启</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import {
  Button,
  Col,
  Descriptions,
  Form, message,
  Row,
  Select,
} from 'ant-design-vue';
import {CrownTwoTone} from "@ant-design/icons-vue"
import emulatorAPI from "../api/emulator.js";
import roomAPI from "../api/room.js";
import roomMemberAPI from "../api/roomMember.js";
import api from "../api/request.js";

export default {
  props: {
    roomId: String,
  },
  components: {
    AButton: Button,
    AForm: Form,
    ARow: Row,
    ACol: Col,
    CrownTwoTone,
    AFormItem: Form.Item,
    ASelect: Select,
    ADescriptions: Descriptions,
    ADescriptionsItem: Descriptions.Item,
  },
  data() {
    return {
      emulatorOptions: [],
      currentUserIsHost: false,
      emulators: [],
      selectedEmulator: null,
      emulatorGames: [],
      emulatorGameOptions: [],
      selectedGame: null,

      updateForm: {
        roomId: null,
        emulatorId: null,
        gameId: null,
      },

      roomDetail: null,
      userRoomMember: null,
    }
  },
  created() {
    this.updateForm.roomId = this.roomId;
    this.getUserRoomMember();
    this.getRoomDetail().then(this.listEmulator).then(this.listGames);
    addEventListener("emulatorInfoDrawerOpen", _=>{
      this.selectedEmulator = null;
      this.selectedGame = null;
      this.getUserRoomMember();
      this.getRoomDetail().then(this.listEmulator).then(this.listGames);
    });
  },
  methods: {
    onEmulatorSelectChange() {
      this.selectedEmulator = this.emulators.find(item => item.emulatorId === this.updateForm.emulatorId);
      if (this.selectedEmulator) {
        this.listGames();
      }
    },

    listEmulator() {
      return emulatorAPI.listEmulator().then(resp => {
        const data = resp.data;
        this.emulators = data;
        this.emulatorOptions = data.map(item => {
          return {
            label: item["emulatorName"],
            value: item["emulatorId"]
          }
        });
        if (this.updateForm.emulatorId !== "0") {
          this.selectedEmulator = this.emulators.find(item => item.emulatorId === this.updateForm.emulatorId);
          if (this.selectedEmulator) return;
        }
        this.updateForm.emulatorId = this.emulators[0]['emulatorId'];
        this.selectedEmulator = this.emulators[0];
      });
    },

    listGames() {
      return emulatorAPI.listGame(this.selectedEmulator.emulatorType).then(resp => {
        const data = resp.data;
        this.emulatorGames = data;
        this.emulatorGameOptions = data.map(item => {
          return {
            label: item["gameName"],
            value: item["gameId"]
          }
        });
        if (this.updateForm.gameId !== "0") {
          this.selectedGame = this.emulatorGames.find(item => item.gameId === this.updateForm.gameId);
          if (this.selectedGame) return;
        }
        this.updateForm.gameId = this.emulatorGames[0]['gameId'];
        this.selectedGame = this.emulatorGames[0];
      });
    },

    getRoomDetail() {
      return roomAPI.getRoomDetail(this.roomId).then(resp=>{
        this.roomDetail = resp.data;
        this.updateForm.emulatorId = this.roomDetail.emulatorId;
        this.updateForm.gameId = this.roomDetail.gameId;
      });
    },

    getUserRoomMember() {
      roomMemberAPI.getUserRoomMember(this.roomId).then(resp => {
        this.userRoomMember = resp.data;
        this.currentUserIsHost = this.userRoomMember["role"] === 1;
      })
    },

    restart() {
      if (this.updateForm.emulatorId === null || this.updateForm.gameId === null) {
        message.warn("请先选择模拟器和游戏");
        return;
      }
      api.post("/room-instance/restart", this.updateForm).then(_=>{
        message.success("重启成功");
      }).catch(resp=>{
        message.error(resp.message);
      })
    }
  }
}
</script>