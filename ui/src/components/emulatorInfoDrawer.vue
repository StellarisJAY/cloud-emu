<template>
  <div>
    <a-form layout="horizontal" labelAlign="right">
      <a-form-item label="模拟器">
        <a-select :options="emulatorOptions" v-model:value="updateForm.emulatorId"
                  @change="onEmulatorSelectChange" :disabled="!currentUserIsHost"></a-select>
      </a-form-item>
      <a-descriptions :column="1" v-if="selectedEmulator" title="模拟器详情">
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
        <a-button type="primary" @click="restart" :disabled="!currentUserIsHost">重启</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import {
  Button,
  Col,
  Descriptions,
  Form,
  Row,
  Select,
} from 'ant-design-vue';
import {CrownTwoTone} from "@ant-design/icons-vue"
import emulatorAPI from "../api/emulator.js";
import roomAPI from "../api/room.js";
import roomMemberAPI from "../api/roomMember.js";

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
        emulatorId: null,
        gameId: null,
      },

      roomDetail: null,
      userRoomMember: null,
    }
  },
  created() {
    this.getUserRoomMember();
    this.getRoomDetail();
    this.listEmulator();
    addEventListener("emulatorInfoDrawerOpen", _=>{
      this.selectedEmulator = null;
      this.selectedGame = null;
      this.getUserRoomMember();
      this.listEmulator();
      this.getRoomDetail();
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
      emulatorAPI.listEmulator().then(resp => {
        const data = resp.data;
        this.emulators = data;
        this.emulatorOptions = data.map(item => {
          return {
            label: item["emulatorName"],
            value: item["emulatorId"]
          }
        });
        this.emulatorOptions.splice(0, 0, {label: "无", value: "0"});
        this.listGames();
      });
    },

    listGames() {
      if (!this.selectedEmulator) {
        this.emulatorGameOptions = [{label: "无", value: "0"}];
        return;
      }
      emulatorAPI.listGame(this.selectedEmulator.emulatorId).then(resp => {
        const data = resp.data;
        this.emulatorGames = data;
        this.emulatorGameOptions = data.map(item => {
          return {
            label: item["gameName"],
            value: item["gameId"]
          }
        });
        this.emulatorGameOptions.splice(0, 0, {label: "无", value: "0"});
      })
    },

    getRoomDetail() {
      roomAPI.getRoomDetail(this.roomId).then(resp=>{
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
      console.log(this.updateForm);
    }
  }
}
</script>