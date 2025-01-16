<script setup>
</script>

<template>
  <!--视频、工具栏-->
  <div id="center-container">
    <div id="center-content">
      <a-row style="height: 15%">
        <a-col :span="6">
          <a-button class="toolbar-button" type="primary" :disabled="saveBtnDisabled" @click="saveGame"
                    style="width: 90%">保存</a-button>
        </a-col>
        <a-col :span="12">
          <a-button class="toolbar-button" type="primary" @click="openRoomMemberDrawer" style="width: 90%">房间设置</a-button>
        </a-col>
        <a-col :span="6">
          <a-button class="toolbar-button" type="primary" :disabled="loadBtnDisabled" @click="openSavedGamesDrawer"
                    style="width: 90%">读档</a-button>
        </a-col>
      </a-row>
      <a-row style="height: 70%">
        <div id="video-mask">
          <video id="video"></video>
        </div>
      </a-row>
      <a-row style="height: 15%">
        <a-col :span="6">
          <a-button class="toolbar-button" type="primary" @click="connect"
                    :disabled="connectBtnDisabled">连接</a-button>
        </a-col>
        <a-col :span="12">
          <a-button class="toolbar-button" type="primary" :disabled="emulatorBtnDisabled"
                    @click="openEmulatorInfoDrawer">模拟器/游戏</a-button>
        </a-col>
        <a-col :span="6">
          <a-button class="toolbar-button" type="primary" @click="openSettingDrawer">设置</a-button>
        </a-col>
      </a-row>
    </div>
  </div>

  <!--房间详情-->
  <a-drawer v-model:open="membersDrawerOpen" placement="right" size="default" title="房间详情">
    <RoomInfoDrawer :member-self="memberSelf" :rtc-session="rtcSession" :full-room-info="fullRoomInfo"
                    :room-id="roomId"></RoomInfoDrawer>
  </a-drawer>
  <!--存档列表-->
  <a-drawer size="default" title="保存游戏" placement="right" v-model:open="savedGameOpen">
    <SaveList :room-id="roomId"></SaveList>
  </a-drawer>
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
  <!--模拟器，游戏选项-->
  <a-drawer v-model:open="emulatorInfoDrawerOpen" placement="right" title="模拟器/游戏">
    <emulator-info-drawer :room-id="roomId"></emulator-info-drawer>
  </a-drawer>

  <div id="left-control-buttons">
    <a-row style="height: 33.33%;">
      <a-col :span="8" id="slot-l1"></a-col>
      <a-col :span="8" id="slot-l2"></a-col>
      <a-col :span="8" id="slot-l3"></a-col>
    </a-row>
    <a-row style="height: 33.33%">
      <a-col :span="8" id="slot-l4"></a-col>
      <a-col :span="8" id="slot-l5"></a-col>
      <a-col :span="8" id="slot-l6"></a-col>
    </a-row>
    <a-row style="height: 33.33%">
      <a-col :span="8" id="slot-l7"></a-col>
      <a-col :span="8" id="slot-l8"></a-col>
      <a-col :span="8" id="slot-l9"></a-col>
    </a-row>
  </div>

  <div id="right-control-buttons">
    <a-row style="height: 33.33%;">
      <a-col :span="8" id="slot-r1"></a-col>
      <a-col :span="8" id="slot-r2"></a-col>
      <a-col :span="8" id="slot-r3"></a-col>
    </a-row>
    <a-row style="height: 33.33%">
      <a-col :span="8" id="slot-r4"></a-col>
      <a-col :span="8" id="slot-r5"></a-col>
      <a-col :span="8" id="slot-r6"></a-col>
    </a-row>
    <a-row style="height: 33.33%">
      <a-col :span="8" id="slot-r7"></a-col>
      <a-col :span="8" id="slot-r8"></a-col>
      <a-col :span="8" id="slot-r9"></a-col>
    </a-row>
  </div>

  <!--实时状态参数-->
  <p id="stats" v-if="configs.showStats">RTT:{{ stats.rtt }}ms FPS:{{ stats.fps }} D:{{formatBytes(stats.bytesPerSecond)}}/s</p>
</template>

<script>
import api from "../api/request.js";
import globalConfigs from "../api/const.js";
import { Row, Col } from "ant-design-vue";
import { Card, Button, Drawer, Select,Switch, Slider } from "ant-design-vue";
import { message } from "ant-design-vue";
import { Form, FormItem, Modal, Input } from "ant-design-vue";
import { ArrowUpOutlined, ArrowDownOutlined, ArrowLeftOutlined, ArrowRightOutlined } from "@ant-design/icons-vue"
import { Tour } from "ant-design-vue";
import RoomInfoDrawer from "../components/roomInfoDrawer.vue";
import SaveList from "../components/saveList.vue";
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
    ASlider: Slider,
    EmulatorInfoDrawer,
  },
  data() {
    return {
      membersDrawerOpen: false,
      settingDrawerOpen: false,
      chatModalOpen: false,
      emulatorInfoDrawerOpen: false,
      savedGameOpen: false,

      roomId: null,
      memberSelf: {},
      fullRoomInfo: {},
      mobileDevice: false,

      rtcSession: {},
      connectToken: null,
      iceCandidates: [],

      connectBtnDisabled: false,
      saveBtnDisabled: true,
      loadBtnDisabled: true,
      chatBtnDisabled: true,
      emulatorBtnDisabled: true,
      graphicOptionsDisabled: true,
      emulatorSpeedSliderDisabled: true,

      pingInterval: 0,
      stats: {
        rtt: 0,
        fps: 0,
        bytesReceived: 0,
        bytesPerSecond: 0,
      },

      configs: {
        showStats: false,
      },
      graphicOptions: {
        highResOpen: false,
        reverseColor: false,
        grayscale: false,
      },
      emulatorSpeedRate: 1.0,
      allowedEmulatorSpeedRates: {
        0.5: "0.5x",
        0.75: "0.75x",
        1.0: "1.0x",
        1.25: "1.25x",
        1.5: "1.5x",
        2.0: "2.0x"
      },

      videoHidden: true,
    }
  },
  created() {
    this.mobileDevice = platform.isMobile();

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
  mounted() {
    if (platform.isPortraitOrientation()) {
      this.initPortraitLayout();
      message.info("请使用横屏全屏来获取最佳游戏体验");
    }else {
      this.initLandscapeLayout();
    }
    window.onresize = this.onWindowResize;
  },
  unmounted() {
    if (this.rtcSession && this.rtcSession.pc) {
      this.rtcSession.pc.close();
    }
  },
  methods: {
    onWindowResize() {
      if (platform.isPortraitOrientation()) {
        this.initPortraitLayout();
      }else {
        this.initLandscapeLayout();
      }
    },
    initPortraitLayout() {
      const center = document.getElementById("center-container");
      const lButtons = document.getElementById("left-control-buttons");
      const rButtons = document.getElementById("right-control-buttons");
      const centerContent = document.getElementById("center-content");
      const video = document.getElementById("video");

      center.style.width = "100%";
      center.style.height = "100%";
      center.style.position = "absolute";
      center.style.left = "0%";

      lButtons.style.position = "absolute";
      lButtons.style.width = "50%";
      lButtons.style.left = "0";
      lButtons.style.height = "50%"
      lButtons.style.top = "50%";

      rButtons.style.position = "absolute";
      rButtons.style.width = "50%";
      rButtons.style.left = "50%";
      rButtons.style.height = "50%"
      rButtons.style.top = "50%";

      centerContent.style.height = "50%";
      video.style.height = "100%";
      video.style.width = "100%";
    },

    initLandscapeLayout() {
      const center = document.getElementById("center-container");
      const lButtons = document.getElementById("left-control-buttons");
      const rButtons = document.getElementById("right-control-buttons");
      const centerContent = document.getElementById("center-content");
      const video = document.getElementById("video");

      center.style.width = "40%";
      center.style.height = "100%";
      center.style.position = "absolute";
      center.style.left = "30%";

      lButtons.style.position = "absolute";
      lButtons.style.width = "30%";
      lButtons.style.left = "0";
      lButtons.style.height = "100%"
      lButtons.style.top = "0";

      rButtons.style.position = "absolute";
      rButtons.style.width = "30%";
      rButtons.style.left = "70%";
      rButtons.style.height = "100%"
      rButtons.style.top = "0%";

      centerContent.style.height = "100%";
      video.style.width = "100%";
      video.style.height = "100%";
    },

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
        this.initScreenButtons(this.roomDetail["emulatorType"]);
        this.setKeyboardControls(this.roomDetail["emulatorType"], false);
      });
      dispatchEvent(new Event("webrtc-connected"));
      this.saveBtnDisabled = false;
      this.loadBtnDisabled = false;
      this.emulatorBtnDisabled = false;
    },
    onDisconnected() {
      message.warn("连接断开");
      dispatchEvent(new Event("webrtc-disconnected"));
      this.destroyScreenButtons();
      this.rtcSession.pc.close();
      this.saveBtnDisabled = true;
      this.loadBtnDisabled = true;
      this.chatBtnDisabled = true;
      this.graphicOptionsDisabled = true;
      this.emulatorBtnDisabled = true;
      this.connectBtnDisabled = false;
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
        case MessageRestart:
          message.info("模拟器重启");
          this.destroyScreenButtons();
          this.initScreenButtons(msgObj["data"]["EmulatorType"]);
          this.setKeyboardControls(msgObj["data"]["EmulatorType"], false);
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

    initScreenButtons: function(emulatorType) {
      // 初始化时使用默认布局
      const layout = globalConfigs.defaultButtonLayouts[emulatorType];
      if (!layout) return;
      layout.forEach(item => {
        let button = document.createElement("button");
        button.type = "button";
        button.classList.add("control-button");
        button.innerText = item.label;
        button.style.width = "100%";
        button.style.height = "100%";
        button.style.backgroundColor = "#f8f3d4"
        button.style.position="relative"
        button.addEventListener("mousedown", _=>this.sendAction(item["code"], MessageGameButtonPressed));
        button.addEventListener("mouseup", _=>this.sendAction(item["code"], MessageGameButtonReleased));
        button.addEventListener("touchstart", _=>this.sendAction(item["code"], MessageGameButtonPressed));
        button.addEventListener("touchend", _=>this.sendAction(item["code"], MessageGameButtonReleased));
        const id = "slot-"+item["slot"];
        document.getElementById(id).appendChild(button);
      });
    },

    setKeyboardControls: function(emulatorType, reset) {
      if (reset) {
        window.onkeyup = _=>{};
        window.onkeydown = _=>{};
        return;
      }
      const binding = globalConfigs.defaultKeyboardBindings[emulatorType];
      if (!binding) return;
      window.onkeydown = (ev)=>{
        binding.forEach(item=>{
          if (item.keyCode === ev.code) {
            this.sendAction(item["button"], MessageGameButtonPressed)
          }
        });
      };
      window.onkeyup = (ev)=>{
        binding.forEach(item=>{
          if (item.keyCode === ev.code) {
            this.sendAction(item["button"], MessageGameButtonReleased)
          }
        });
      };
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
  left: 0;
  right: 0;
  margin: auto;
}

#video-mask {
  display:flex;
  background-color: black;
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
}

#center-container {
  background-color: #f8f3d4;
  position: absolute;
  left: 30%;
  width: 40%;
  top: 0;
  bottom: 0;
  padding: 1%;
}

.toolbar-button {
  background-color: #9c253d;
  width: 90%;
}

.toolbar-button:hover {
  background-color: #811f33;
}

#stats {
  position: absolute;
  right: 0;
  top: 0;
}

#left-control-buttons {
  position: absolute;
  left: 0;
  top: 0;
  width: 30%;
  height: 100%;
}

#right-control-buttons {
  position: absolute;
  right: 0;
  top: 0;
  width: 30%;
  height: 100%;
}
</style>