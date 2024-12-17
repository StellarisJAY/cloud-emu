<template>
  <div>
    <a-form>
      <a-form-item label="访问权限">
        <a-select :options="joinTypeOptions" v-model:value="updateRoomForm.joinType"></a-select>
      </a-form-item>
      <a-form-item label="模拟器">
        <a-select :options="emulatorOptions" v-model:value="updateRoomForm.emulatorId"
          @change="onEmulatorSelectChange"></a-select>
      </a-form-item>
      <a-descriptions :column="1" v-if="selectedEmulator" title="模拟器详情">
        <a-descriptions-item label="描述">{{ selectedEmulator.description }}</a-descriptions-item>
        <a-descriptions-item label="提供者">{{ selectedEmulator.provider }}</a-descriptions-item>
        <a-descriptions-item label="支持存档">{{ selectedEmulator.supportSave ? "是" : "否" }}</a-descriptions-item>
        <a-descriptions-item label="支持画面设置">{{ selectedEmulator.supportGraphicSetting ? "是" : "否"
          }}</a-descriptions-item>
      </a-descriptions>
      <a-form-item label="游戏">
        <a-select :options="emulatorGameOptions" v-model:value="updateRoomForm.gameId"></a-select>
      </a-form-item>
    </a-form>
    <a-list item-layout="vertical" :data-source="members">
      <template #renderItem="{ item }">
        <a-list-item>
          <a-row>
            <a-col :span="8">
              <CrownTwoTone v-if="item.role === 1" />{{ item.nickName }}
            </a-col>
            <a-col :span="10">
            </a-col>
          </a-row>
        </a-list-item>
      </template>
    </a-list>
  </div>
</template>

<script>
import { List, Button, Checkbox, Drawer, Form, Input, Switch, Radio, RadioGroup, Select, Descriptions } from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import { message } from "ant-design-vue";
import api from "../api/request.js";
import { CrownTwoTone } from "@ant-design/icons-vue"
import constants from "../api/const.js";
import emulatorAPI from "../api/emulator.js";

export default {
  props: {
    memberSelf: Object,
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
    ASwitch: Switch,
    ACheckbox: Checkbox,
    AList: List,
    AListItem: List.Item,
    ARadio: Radio,
    ARadioGroup: RadioGroup,
    CrownTwoTone,
    AInputPassword: Input.Password,
    AFormItem: Form.Item,
    ASelect: Select,
    ADescriptions: Descriptions,
    ADescriptionsItem: Descriptions.Item,
  },
  data() {
    return {
      members: [],
      RoleNameHost: "Host",
      RoleNamePlayer: "Player",
      RoleNameObserver: "Observer",
      privacySwitchDisabled: false,

      joinTypeOptions: [],
      emulatorOptions: [],
      updateRoomForm: {
        joinType: 1,
        emulatorId: '1',
        gameId: null
      },
      emulators: [],
      selectedEmulator: null,
      emulatorGames: [],
      emulatorGameOptions: [],
      selectedGame: null,
    }
  },
  created() {
    this.joinTypeOptions = constants.getEnumOptions("roomJoinTypeEnum");
    this.listRoomMembers();
    let _this = this;
    emulatorAPI.listEmulator().then(resp => {
      const data = resp.data;
      _this.emulators = data;
      _this.emulatorOptions = data.map(item => {
        return {
          label: item.emulatorName,
          value: item.emulatorId
        }
      });
      _this.selectedEmulator = data[0];
      _this.updateRoomForm.emulatorId = data[0].emulatorId;
      _this.listGames();
    });
    addEventListener("memberDrawerOpen", _ => this.listRoomMembers());
  },
  methods: {
    listRoomMembers: async function () {
      const resp = await api.get("/room-member?roomId=" + this.roomId);
      const data = resp.data;
      this.members = data;
    },
    onEmulatorSelectChange() {
      this.selectedEmulator = this.emulators.find(item => item.emulatorId === this.updateRoomForm.emulatorId);
      if (this.selectedEmulator) {
        this.listGames();
      }
    },
    listGames() {
      let _this = this;
      emulatorAPI.listGame(this.selectedEmulator.emulatorId).then(resp => {
        const data = resp.data;
        _this.emulatorGames = data;
        _this.emulatorGameOptions = data.map(item => {
          return {
            label: item.gameName,
            value: item.gameId
          }
        });
        _this.selectedGame = data[0];
        _this.updateRoomForm.gameId = data[0].gameId;
      })
    }
  }
}
</script>