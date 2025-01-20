<template>
    <a-layout id="layout">
        <a-layout-header style="position: fixed; z-index: 1; width: 100%;">
            <a-menu theme="dark" mode="horizontal" v-model:selectedKeys="headerSelectedKeys">
                <a-menu-item key="1">我的房间</a-menu-item>
                <a-menu-item key="2">加入房间</a-menu-item>
                <a-menu-item key="3">通知</a-menu-item>
                <a-menu-item key="4">个人</a-menu-item>
                <a-menu-item key="6" v-if="isAdmin">游戏列表</a-menu-item>
                <a-menu-item key="7" v-if="isAdmin">模拟器</a-menu-item>
              <a-menu-item key="8" v-if="isAdmin">宏设置</a-menu-item>
                <a-menu-item key="5" @click="logout">注销</a-menu-item>
            </a-menu>
        </a-layout-header>
        <a-row style="margin-top: 64px;">
            <a-col :xs="{ offset: 2, span: 20 }" :sm="{ offset: 2, span: 20 }" :md="{ offset: 4, span: 16 }"
                :lg="{ offset: 4, span: 16 }">
                <RoomList v-if="headerSelectedKeys[0] === '1'" :joined="true" />
                <RoomList v-else-if="headerSelectedKeys[0] === '2'" :joined="false"></RoomList>
                <NotificationList v-else-if="headerSelectedKeys[0] === '3'"></NotificationList>
                <UserInfo v-else-if="headerSelectedKeys[0] === '4'"/>
                <AdminGameList v-else-if="headerSelectedKeys[0] === '6'"></AdminGameList>
                <AdminEmulatorList v-else-if="headerSelectedKeys[0] === '7'"/>
                <MacroSettings v-else-if="headerSelectedKeys[0] === '8'"/>
            </a-col>
        </a-row>
    </a-layout>
</template>

<script>
import { Layout, Menu, Card } from 'ant-design-vue';
import { Row, Col } from "ant-design-vue";
import RoomList from "../components/roomList.vue";
import router from "../router/index.js";
import tokenStorage from "../api/token.js";
import NotificationList from "../components/notificationList.vue";
import UserInfo from "../components/userInfo.vue";
import AdminInfo from "../components/adminGameList.vue";
import userAPI from "../api/user.js";
import AdminGameList from "../components/adminGameList.vue";
import AdminEmulatorList from "../components/adminEmulatorList.vue";
import MacroSettings from "../components/macroSettings.vue";

export default {
    components: {
      AdminGameList,
        ALayout: Layout,
        ALayoutHeader: Layout.Header,
        AMenu: Menu,
        AMenuItem: Menu.Item,
        ALayoutContent: Layout.Content,
        ARow: Row,
        ACol: Col,
        ACard: Card,
        RoomList: RoomList,
        NotificationList: NotificationList,
        UserInfo: UserInfo,
        AdminInfo: AdminInfo,
        AdminEmulatorList: AdminEmulatorList,
        MacroSettings: MacroSettings,
    },
    data() {
        return {
            headerSelectedKeys: ['1'],
            isAdmin: false,
        }
    },
    created() {
      userAPI.getLoginUserDetail().then(res=>{
          this.isAdmin = res.data.role === 2;
      });
    },
  methods: {
        logout: function () {
            tokenStorage.delToken();
            router.push("/login");
        },
    },
}
</script>

<style>
#layout {
  background-color: #e1e1e1;
}
</style>