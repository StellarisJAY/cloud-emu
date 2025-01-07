<script setup>
import {ref} from "vue";

const refConnBtn = ref(null)
const refSelector = ref(null)
const refRestart = ref(null)
const refSaveBtn = ref(null)
const refLoadBtn = ref(null)
const refRoomBtn = ref(null)
const refSettingBtn = ref(null)
const tourSteps = [
  {
    title: "选择游戏",
    description: "点击此处选择需要加载的游戏",
    target: () => refSelector.value && refSelector.value.$el,
  },
  {
    title: "连接",
    description: "点击按钮连接到模拟器",
    target: () => refConnBtn.value && refConnBtn.value.$el,
  },
  {
    title: "重启",
    description: "重启可以用于切换游戏，但会清除当前游戏进度，如有需要请先保存当前游戏。",
    target: () => refRestart.value && refRestart.value.$el,
  },
  {
    title: "保存游戏",
    description: "点击此处保存游戏进度，请注意存档数量上限。",
    target: () => refSaveBtn.value && refSaveBtn.value.$el,
  },
  {
    title: "读取存档",
    description: "显示存档列表，跨游戏读取存档会重启模拟器，如有需要请先保存当前游戏。",
    target: () => refLoadBtn.value && refLoadBtn.value.$el,
  },
  {
    title: "房间管理",
    description: "点击此处弹出房间面板，房主可通过此面板修改房间信息以及玩家权限。",
    target: () => refRoomBtn.value && refRoomBtn.value.$el,
  },
  {
    title: "游戏设置",
    description: "点击此处弹出设置面板，可设置游戏图像，键盘按键绑定。",
    target: ()=> refSettingBtn.value && refSettingBtn.value.$el,
  }
]
</script>

<template>
  <a-row style="height: 100vh; background-color: #00b8a9">
    <!--左侧控制按钮-->
    <a-col :span="7">
      <a-row style="height: 30%; margin-top: 10%">
        <a-col :span="8" id="slot-l1"></a-col>
        <a-col :span="8" id="slot-l2"></a-col>
        <a-col :span="8" id="slot-l3"></a-col>
      </a-row>
      <a-row style="height: 30%">
        <a-col :span="8" id="slot-l4"></a-col>
        <a-col :span="8" id="slot-l5"></a-col>
        <a-col :span="8" id="slot-l6"></a-col>
      </a-row>
      <a-row style="height: 30%">
        <a-col :span="8" id="slot-l7"></a-col>
        <a-col :span="8" id="slot-l8"></a-col>
        <a-col :span="8" id="slot-l9"></a-col>
      </a-row>
    </a-col>
    <!--视频、工具栏-->
    <a-col :span="10">
      <a-card id="center-container">
        <a-row>
          <a-col :span="6">
            <a-button ref="refSaveBtn" class="toolbar-button" type="primary" :disabled="saveBtnDisabled" @click="saveGame"
              style="width: 90%">保存</a-button>
          </a-col>
          <a-col :span="6">
            <a-button ref="refLoadBtn" class="toolbar-button" type="primary" :disabled="loadBtnDisabled" @click="openSavedGamesDrawer"
              style="width: 90%">读档</a-button>
          </a-col>
          <a-col :span="6">
            <a-button class="toolbar-button" type="primary" :disabled="chatBtnDisabled" @click="_ => { setChatModal(true) }"
              style="width: 90%">聊天</a-button>
          </a-col>
          <a-col :span="6">
            <a-button ref="refRoomBtn" class="toolbar-button" type="primary" @click="openRoomMemberDrawer" style="width: 90%">房间</a-button>
          </a-col>
        </a-row>
        <a-row style="height:80%;margin-bottom:20px;margin-top:20px;">
          <video id="video" playsinline></video>
        </a-row>
        <a-row>
          <a-col :span="6">
            <a-button ref="refConnBtn" class="toolbar-button" type="primary" @click="connect"
              :disabled="connectBtnDisabled">连接</a-button>
          </a-col>
          <a-col :span="12">
            <a-button ref="refRestart" class="toolbar-button" type="primary"
              @click="openEmulatorInfoDrawer">模拟器/游戏</a-button>
          </a-col>
          <a-col :span="6">
            <a-button ref="refSettingBtn" class="toolbar-button" type="primary" @click="openSettingDrawer">设置</a-button>
          </a-col>
        </a-row>
      </a-card>
    </a-col>
    <!--右侧控制按钮-->
    <a-col :span="7">
      <a-row style="height: 30%; margin-top: 10%">
        <a-col :span="8" id="slot-r1"></a-col>
        <a-col :span="8" id="slot-r2"></a-col>
        <a-col :span="8" id="slot-r3"></a-col>
      </a-row>
      <a-row style="height: 30%">
        <a-col :span="8" id="slot-r4"></a-col>
        <a-col :span="8" id="slot-r5"></a-col>
        <a-col :span="8" id="slot-r6"></a-col>
      </a-row>
      <a-row style="height: 30%">
        <a-col :span="8" id="slot-r7"></a-col>
        <a-col :span="8" id="slot-r8"></a-col>
        <a-col :span="8" id="slot-r9"></a-col>
      </a-row>
    </a-col>
    <!--房间详情-->
    <a-drawer v-model:open="membersDrawerOpen" placement="right" size="default" title="房间详情">
      <RoomInfoDrawer :member-self="memberSelf" :rtc-session="rtcSession" :full-room-info="fullRoomInfo"
        :room-id="roomId"></RoomInfoDrawer>
    </a-drawer>
    <!--存档列表-->
    <a-drawer size="default" title="保存游戏" placement="right" v-model:open="savedGameOpen">
      <SaveList :room-id="roomId"></SaveList>
    </a-drawer>
    <!--聊天窗口-->
    <a-modal title="聊天" v-model:open="chatModalOpen" @cancel="_ => { setChatModal(false) }">
      <template #footer>
        <a-button @click="_ => { setChatModal(false) }">取消</a-button>
        <a-button type="primary" @click="sendChatMessage">发送</a-button>
      </template>
      <a-input placeholder="请输入消息..." v-model:value="chatMessage"></a-input>
    </a-modal>
    <!--设置列表-->
    <a-drawer v-model:open="settingDrawerOpen" placement="right" title="设置" size="default">
      <a-form>
        <a-form-item label="显示状态数据">
          <a-switch v-model:checked="configs.showStats"></a-switch>
        </a-form-item>
      </a-form>
      <a-form>
        <a-form-item label="高分辨率">
          <a-switch v-model:checked="graphicOptions.highResOpen" :disabled="graphicOptionsDisabled" @change="updateGraphicOptions"></a-switch>
        </a-form-item>
        <a-form-item label="反色">
          <a-switch v-model:checked="graphicOptions.reverseColor" :disabled="graphicOptionsDisabled" @change="updateGraphicOptions"></a-switch>
        </a-form-item>
        <a-form-item label="黑白">
          <a-switch v-model:checked="graphicOptions.grayscale" :disabled="graphicOptionsDisabled" @change="updateGraphicOptions"></a-switch>
        </a-form-item>
        <a-form-item label="模拟器速度">
          <a-slider v-model:value="emulatorSpeedRate" :min="0.5" :max="2.0" :marks="allowedEmulatorSpeedRates" :step="null" @afterChange="setEmulatorSpeed" :disabled="emulatorSpeedSliderDisabled"></a-slider>
        </a-form-item>
      </a-form>
    </a-drawer>
    <a-tour :steps="tourSteps" :open="tourOpen" @close="_ => { tourOpen = false }"></a-tour>
    <!--模拟器，游戏选项-->
    <a-drawer v-model:open="emulatorInfoDrawerOpen" placement="right" title="模拟器/游戏">
      <emulator-info-drawer :room-id="roomId"></emulator-info-drawer>
    </a-drawer>
  </a-row>
  <!--实时状态参数-->
  <p id="stats" v-if="configs.showStats">RTT:{{ stats.rtt }}ms FPS:{{ stats.fps }} D:{{formatBytes(stats.bytesPerSecond)}}/s</p>
</template>

<script>
import api from "../api/request.js";
import globalConfigs from "../api/const.js";
import { Row, Col } from "ant-design-vue";
import { Card, Button, Drawer, Select,Switch, notification, Slider } from "ant-design-vue";
import { message } from "ant-design-vue";
import { Form, FormItem, Modal, Input } from "ant-design-vue";
import { ArrowUpOutlined, ArrowDownOutlined, ArrowLeftOutlined, ArrowRightOutlined } from "@ant-design/icons-vue"
import { Tour } from "ant-design-vue";
import RoomInfoDrawer from "../components/roomInfoDrawer.vue";
import SaveList from "../components/saveList.vue";
import KeyboardSetting from "../components/keyboardSetting.vue";
import platform from "../util/platform.js";
import EmulatorInfoDrawer from "../components/emulatorInfoDrawer.vue";
import roomMemberAPI from "../api/roomMember.js";
import router from "../router/index.js";
import roomAPI from "../api/room.js";

const MessageGameButtonPressed = 0
const MessageGameButtonReleased = 1
const MessageChat = 2;
const MessagePing = 14;
const MessageRestart = 17;

export default {
  components: {
    ARow: Row,
    ACol: Col,
    ACard: Card,
    AButton: Button,
    ADrawer: Drawer,
    ArrowUpOutlined,
    ArrowDownOutlined,
    ArrowLeftOutlined,
    ArrowRightOutlined,
    ASelect: Select,
    ASwitch: Switch,
    AForm: Form,
    AFormItem: FormItem,
    ATour: Tour,
    AInput: Input,
    AModal: Modal,
    AInputPassword: Input.Password,
    RoomInfoDrawer,
    SaveList,
    KeyboardSetting,
    ASlider: Slider,
    EmulatorInfoDrawer,
  },
  data() {
    return {
      membersDrawerOpen: false,
      memberSelf: {},
      rtcSession: {},
      connectBtnDisabled: false,
      saveBtnDisabled: true,
      loadBtnDisabled: true,
      restartBtnDisabled: true,
      chatBtnDisabled: true,
      selectedGame: "",
      configs: {
        controlButtonMapping: {
          "button-up": "Up",
          "button-down": "Down",
          "button-left": "Left",
          "button-right": "Right",
          "button-a": "A",
          "button-b": "B",
          "button-select": "Select",
          "button-start": "Start",
        },
        showStats: false,
      },
      savedGameOpen: false,
      fullRoomInfo: {},
      chatModalOpen: false,
      chatMessage: "",
      pingInterval: 0,
      iceCandidates: [],
      settingDrawerOpen: false,
      stats: {
        rtt: 0,
        fps: 0,
        bytesReceived: 0,
        bytesPerSecond: 0,
      },
      joinRoomFormState: {
        id: 0,
        password: "",
      },
      joinRoomModalOpen: false,
      tourOpen: false,
      graphicOptions: {
        highResOpen: false,
        reverseColor: false,
        grayscale: false,
      },
      graphicOptionsDisabled: true,
      mobileDevice: false,
      emulatorSpeedRate: 1.0,
      allowedEmulatorSpeedRates: {
        0.5: "0.5x",
        0.75: "0.75x",
        1.0: "1.0x",
        1.25: "1.25x",
        1.5: "1.5x",
        2.0: "2.0x"
      },
      emulatorSpeedSliderDisabled: true,
      emulatorInfoDrawerOpen: false,

      connectToken: null,
    }
  },
  created() {
    this.mobileDevice = platform.isMobile();
    if (platform.isPortraitOrientation()) {
      message.info("请使用横屏全屏来获取最佳游戏体验");
    }
    this.roomId = this.$route["params"]["roomId"];
    roomMemberAPI.getUserRoomMember(this.roomId).then(res => {
      this.memberSelf = res.data;
    }).catch(resp=>{
      if (resp["code"] && resp["code"] === 404) {
        router.push("/join-room/" + this.roomId);
      } else {
        message.error("无法进入房间");
      }
    });
  },
  unmounted() {
    if (this.rtcSession && this.rtcSession.pc) {
      this.rtcSession.pc.close();
    }
    this.setKeyboardControl(false);
  },
  methods: {
    openRoomMemberDrawer() {
      this.membersDrawerOpen = true;
      dispatchEvent(new Event("memberDrawerOpen"));
    },

    openEmulatorInfoDrawer() {
      this.emulatorInfoDrawerOpen = true;
      dispatchEvent(new Event("emulatorInfoDrawerOpen"));
    },

    connect() {
      this.connectBtnDisabled = true
      this.openConnection();
    },

    openConnection: async function () {
      const roomId = this.roomId;
      try {
        const response = await api.get("/room-instance", {"roomId": this.roomId});
        this.connectToken = response.data["accessToken"];
        const resp = await api.post("/room-instance/connect", {
          "roomId": roomId,
          "token": this.connectToken,
        });
        await this.createWebRTCPeerConnection(resp.data);
      } catch (errResp) {
        message.warn("连接失败，请重试");
        this.connectBtnDisabled = false;
      }
    },
    createWebRTCPeerConnection: async function (data) {
      const pc = new RTCPeerConnection({
        iceServers: [
          {
            urls: globalConfigs.StunServer,
          },
          {
            urls: globalConfigs.TurnServer.Host,
            username: globalConfigs.TurnServer.Username,
            credential: globalConfigs.TurnServer.Password,
          }
        ],
        iceTransportPolicy: "all",
      });
      // on remote track
      pc.ontrack = ev => {
        console.log("on track: ", ev);
        if (ev.track.kind === "video") {
          document.getElementById("video").srcObject = ev.streams[0]
          document.getElementById("video").autoplay = true
          document.getElementById("video").controls = true
        }
      };

      // 发送answer之前的candidate，避免远端没有收到answer导致无法这是candidate
      pc.onicecandidate = ev => {
        if (ev.candidate) {
          this.iceCandidates.push(ev.candidate);
        }
      };

      pc.onconnectionstatechange = ev => this.onPeerConnStateChange(ev);

      const rtcSession = {
        roomId: this.roomId,
        pc: pc,
        dataChannel: null,
      }
      await pc.setRemoteDescription({
        type: "offer",
        sdp: data["sdpOffer"],
      });
      const answer = await pc.createAnswer();
      await pc.setLocalDescription(answer);
      try {
        await api.post("/room-instance/sdp-answer", {
          "roomId": this.roomId,
          "token": this.connectToken,
          "sdpAnswer": answer.sdp,
        });
      } catch (errResp) {
        message.warn("连接失败，请重试");
        this.connectBtnDisabled = false;
        return;
      }

      // 发送answer之前的candidate，避免远端没有收到answer导致无法这是candidate
      this.iceCandidates.forEach(candidate => {
        const s = JSON.stringify(candidate);
        api.post("/room-instance/ice-candidate", {
          "roomId": this.roomId,
          "token": this.connectToken,
          "iceCandidate": s,
        }).then(_ => {
          this.getRemoteCandidate();
        });
      });
      // 发送answer之后的candidate直接发送给远端
      pc.onicecandidate = ev => {
        if (ev.candidate) {
          const s = JSON.stringify(ev.candidate);
          console.log(ev.candidate);
          api.post("/room-instance/ice-candidate", {
            "roomId": this.roomId,
            "token": this.connectToken,
            "iceCandidate": s,
          }).then(_ => {
            this.getRemoteCandidate();
          });
        }
      }
      // data channel
      pc.ondatachannel = ev => {
        rtcSession.dataChannel = ev.channel;
        ev.channel.onopen = _ => this.onDataChannelOpen();
        ev.channel.onmessage = msg => this.onDataChannelMsg(msg);
        ev.channel.onclose = _ => this.onDataChannelClose();
      };
      this.rtcSession = rtcSession;
    },

    getRemoteCandidate() {
      api.get("/room-instance/ice-candidate", {"roomId": this.roomId, "token": this.connectToken}).then(resp=>{
        resp.data.forEach(candidate => {
          const c = JSON.parse(candidate);
          if (c && candidate !== "") {
            this.rtcSession.pc.addIceCandidate(c);
          }
        });
      });
    },

    onPeerConnStateChange(_) {
      const pc = this.rtcSession.pc
      console.log("peer conn state: " + pc.connectionState)
      switch (pc.connectionState) {
        case "connected":
          this.onConnected()
          break
        case "disconnected":
          this.onDisconnected()
          break
        default:
          break
      }
    },
    onConnected() {
      message.success("连接成功");
      roomAPI.getRoomDetail(this.roomId).then(resp=>{
        this.roomDetail = resp.data;
        this.initScreenButtons(this.roomDetail["emulatorId"]);
      });
      dispatchEvent(new Event("webrtc-connected"));
      this.restartBtnDisabled = false;
      this.saveBtnDisabled = false;
      this.loadBtnDisabled = false;
    },
    onDisconnected() {
      message.warn("连接断开");
      dispatchEvent(new Event("webrtc-disconnected"));
      this.destroyScreenButtons();
      this.rtcSession.pc.close();
      this.saveBtnDisabled = true;
      this.restartBtnDisabled = true;
      this.loadBtnDisabled = true;
      this.chatBtnDisabled = true;
      this.graphicOptionsDisabled = true;
    },
    sendAction(code, pressed) {
      const msg = JSON.stringify({
        "type": pressed,
        "data": code,
      });
      this.rtcSession.dataChannel.send(msg);
    },

    openSavedGamesDrawer() {
      this.savedGameOpen = true;
      dispatchEvent(new Event("saveListOpen"));
    },
    saveGame() {
      this.saveBtnDisabled = true;
      api.post("/game-save/save", {
        "roomId": this.roomId,
      }).then(_ => {
        message.success("保存成功");
        this.saveBtnDisabled = false;
      }).catch(resp => {
        message.error(resp.message);
        this.saveBtnDisabled = false;
      })
    },

    setChatModal(open) {
      this.setKeyboardControl(!open)
      this.chatModalOpen = open
      if (!open) {
        this.chatMessage = ""
      }
    },
    sendChatMessage() {
      const timestamp = new Date().getTime();
      if (this.rtcSession && this.rtcSession.pc) {
        const pingMsg = {
          "type": MessageChat,
          "timestamp": timestamp,
          "data": this.chatMessage,
        };
        this.rtcSession.dataChannel.send(JSON.stringify(pingMsg));
      }
      this.setChatModal(false);
    },

    setKeyboardControl(enabled) {
      if (enabled) {
        const _this = this;
        let setting;
        if (this.$refs.refKeyboardSettings && this.$refs.refKeyboardSettings.selected) {
          setting = this.$refs.refKeyboardSettings.selected;
        }else {
          setting = globalConfigs.defaultKeyboardSetting;
        }
        window.onkeydown = ev => {
          const button = setting.bindings.find(item => item.buttons[0] === ev.code);
          if (button) {
            _this.sendAction(button.emulatorKey, MessageGameButtonPressed);
          }
        };

        window.onkeyup = ev => {
          const button = setting.bindings.find(item => item.buttons[0] === ev.code);
          if (button) {
            _this.sendAction(button.emulatorKey, MessageGameButtonReleased);
          }
        };
      } else {
        window.onkeydown = _ => { }
        window.onkeyup = _ => { }
      }
    },

    ping() {
      const timestamp = new Date().getTime();
      if (this.rtcSession && this.rtcSession.pc) {
        const pingMsg = {
          "type": MessagePing,
          "timestamp": timestamp,
        };
        this.rtcSession.dataChannel.send(JSON.stringify(pingMsg));
      }
    },

    openSettingDrawer() {
      this.settingDrawerOpen = true;
    },

    onDataChannelMsg(msg) {
      const msgStr = String.fromCharCode.apply(null, new Uint8Array(msg.data));
      const msgObj = JSON.parse(msgStr);
      switch (msgObj.type) {
        case MessagePing:
          this.stats.rtt = new Date().getTime() - msgObj.timestamp;
          break;
        case MessageChat:
          if (!msgObj["from"]) return;
          api.get("/user/" + msgObj["from"]).then(resp => {
            notification.info({
              message: resp["data"]["name"],
              description: msgObj.data,
              placement: "topLeft",
              duration: 1,
            });
          });
          break;
        case MessageRestart:
          message.info("模拟器重启");
          this.destroyScreenButtons();
          this.initScreenButtons(msgObj["data"]["EmulatorId"]);
          break;
        default:
          break
      }
    },

    onDataChannelOpen() {
      this.chatBtnDisabled = false;
      const _this = this;
      this.pingInterval = setInterval(_ => {
        _this.ping();
        _this.collectRTCStats();
      }, 1000);
    },

    onDataChannelClose() {
      if (this.pingInterval) clearInterval(this.pingInterval);
      this.chatBtnDisabled = true;
    },

    collectRTCStats: function () {
      const _this = this;
      if (this.rtcSession && this.rtcSession.pc) {
        const pc = this.rtcSession.pc;
        pc.getStats().then(stats => {
          stats.forEach(report => {
            if (report.type === "inbound-rtp" && report.kind === "video") {
              _this.stats.fps = report.framesPerSecond;
              _this.stats.bytesPerSecond = report.bytesReceived - _this.stats.bytesReceived;
              _this.stats.bytesReceived = report.bytesReceived;
            }
          });
        });
      }
    },

    formatBytes: function (bytes) {
      if(bytes <= 1024) return bytes + "B";
      if(bytes <= 1024*1024) return (bytes>>10) + "KB";
      return (bytes>>20) + "MB";
    },

    updateGraphicOptions: function() {
      const _this = this;
      this.graphicOptionsDisabled = true;
      api.post("/game/graphic", {
        "roomId": this.roomId,
        "highResOpen": this.graphicOptions.highResOpen,
        "reverseColor": this.graphicOptions.reverseColor,
        "grayscale": this.graphicOptions.grayscale,
      }).then(resp => {
        _this.graphicOptionsDisabled = false;
        _this.graphicOptions.highResOpen = resp['highResOpen'];
        _this.graphicOptions.reverseColor = resp['reverseColor'];
        _this.graphicOptions.grayscale = resp['grayscale'];
        message.success("设置成功");
      }).catch(_ => {
        message.error("设置失败");
        _this.graphicOptionsDisabled = false;
      });
    },

    setEmulatorSpeed: function() {
      const _this = this;
      this.emulatorSpeedSliderDisabled = true;
      api.post("/game/speed", {
        "roomId": this.roomId,
        "rate": this.emulatorSpeedRate
      }).then(resp=>{
        _this.emulatorSpeedSliderDisabled = false;
        _this.emulatorSpeedRate = resp["rate"];
        message.success("设置成功");
      }).catch(_=>{
        _this.emulatorSpeedSliderDisabled = false;
        message.error("设置失败");
      });
    },

    initScreenButtons: function(emulatorId) {
      // 初始化时使用默认布局
      const layout = globalConfigs.defaultButtonLayouts[emulatorId];
      layout.forEach(item => {
        let button = document.createElement("button");
        button.type = "button";
        button.classList.add("control-button");
        button.innerText = item.label;
        button.style.width = "100%";
        button.style.height = "100%";
        button.style.backgroundColor = "#f8f3d4"
        button.addEventListener("mousedown", _=>this.sendAction(item["code"], MessageGameButtonPressed));
        button.addEventListener("mouseup", _=>this.sendAction(item["code"], MessageGameButtonReleased));
        button.addEventListener("touchstart", _=>this.sendAction(item["code"], MessageGameButtonPressed));
        button.addEventListener("touchend", _=>this.sendAction(item["code"], MessageGameButtonReleased));
        const id = "slot-"+item["slot"];
        document.getElementById(id).appendChild(button);
      });
    },

    destroyScreenButtons: function() {
      const slots = ["l1","l2","l3","l4","l5","l6","l7","l8","l9","r1","r2","r3","r4","r5","r6","r7","r8","r9"];
      slots.forEach(slot => {
        const element = document.getElementById("slot-"+slot);
        if (element.firstElementChild) element.removeChild(element.firstElementChild);
      });
    },
  }
}
</script>

<style scoped>
#video {
  width: 100%;
  background-color: black;
}

#center-container {
  height: 100vh;
  background-color: #f8f3d4;
}

.toolbar-button {
  background-color: #9c253d;
  width: 90%;
}

.toolbar-button:hover {
  background-color: #811f33;
}

.control-button {
  width: 100%;
  height: 100%;
  background-color: #000;
}

#stats {
  position: absolute;
  right: 0;
  top: 0;
}
</style>