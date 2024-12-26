<script setup>
</script>

<template>
  <a-modal :open="visible" title="选择用户" :closable="false">
    <template #footer>
      <a-button @click="handleCancel">取消</a-button>
      <a-button type="primary" @click="handleOk">确认</a-button>
    </template>
    <a-form>
     <a-row>
       <a-form-item label="用户名">
         <a-input v-model:value="query.userName"></a-input>
       </a-form-item>
       <a-form-item label="昵称">
         <a-input v-model:value="query.nickName"></a-input>
       </a-form-item>
       <a-button @click="listUsers">查询</a-button>
     </a-row>
    </a-form>
    <a-table
        :row-selection="{ selectedRowKeys: selectedRowKeys, onChange: onSelectChange }"
        :columns="columns"
        :data-source="users"
        @change="onPaginationChange"
        :pagination="{pageSize: pageSize, current: page, total: totalPages}"
        row-key="userId"
    />
  </a-modal>
</template>

<script>
import {Row, Col} from "ant-design-vue";
import { Button, Pagination, Modal, Table, Form, Input } from "ant-design-vue";
import userAPI from "../api/user.js";

export default {
  props: {
    visible: Boolean,
  },
  components: {
    ARow: Row,
    ACol: Col,
    AButton: Button,
    APagination: Pagination,
    AModal: Modal,
    ATable: Table,
    AForm: Form,
    AFormItem: Form.Item,
    AInput: Input,
  },
  data() {
    return {
      users: [],
      page: 1,
      pageSize: 10,
      totalPages: 1,
      selectedUsers: [],
      selectedRowKeys: [],
      columns: [
        {title: "用户名", dataIndex: "userName"},
        {title: "昵称", dataIndex: "nickName"},
      ],
      query: {
        userName: "",
        nickName: "",
      }
    }
  },
  created() {
    this.listUsers();
  },
  unmounted() {
  },
  methods: {
    listUsers() {
      userAPI.listUser({
        page: this.page,
        pageSize: this.pageSize,
        userName: this.query.userName,
        nickName: this.query.nickName,
      }).then(resp=>{
        this.users = resp.data;
        this.totalPages = Math.ceil(resp.total / this.pageSize);
      });
    },
    handleOk() {
      const ev = new CustomEvent("userPickerConfirm", {
        detail: {
          "users": this.selectedUsers,
        }
      });
      dispatchEvent(ev);
    },
    handleCancel() {
      dispatchEvent(new Event("userPickerCancel"));
    },
    onSelectChange(ev) {
      this.selectedRowKeys = ev;
      this.selectedUsers = this.users.filter(item=>this.selectedRowKeys.includes(item["userId"]));
    },
    onPaginationChange(pagination) {
      this.page = pagination.current;
      this.pageSize = pagination.pageSize;
      this.listUsers();
    },
  }
}
</script>

<style scoped>
</style>