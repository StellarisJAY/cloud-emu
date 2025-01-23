<template>
  <a-card :bordered="false" class="center-card">
    <a-table :columns="columns" size="small"
             :data-source="emulators" :pagination="false" style="width: 100%">
      <template #bodyCell="{ column, text, record }">
        <template v-if="column.dataIndex === 'supportSave' || column.dataIndex === 'supportGraphicSetting' || column.dataIndex === 'disabled'">
          <a-checkbox v-model:checked="record[column.dataIndex]" @change="updateEmulator(record)"></a-checkbox>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script>
import {Row, Col, Card, message} from "ant-design-vue";
import { Button, Table, Form, Input, Checkbox } from "ant-design-vue";
import emulatorAPI from "../api/emulator.js";
import api from "../api/request.js";

export default {
  components: {
    ARow: Row,
    ACol: Col,
    AButton: Button,
    ACard: Card,
    ATable: Table,
    AFormItem: Form.Item,
    AForm: Form,
    AInput: Input,
    ACheckbox: Checkbox,
  },
  data() {
    return {
      emulators: [],
      columns: [
        {title: "模拟器名称", dataIndex: "emulatorName" },
        {title: "类型", dataIndex: "emulatorType" },
        {title: "允许存档", dataIndex: "supportSave" },
        {title: "允许画面设置", dataIndex: "supportGraphicSetting" },
        {title: "禁用", dataIndex: "disabled" },
      ]
    }
  },
  created() {
    this.listEmulators();
  },
  unmounted() {
  },
  methods: {
    listEmulators() {
      api.get("/emulator", {showDisabled: true}).then(resp=>{
        this.emulators = resp.data;
      });
    },
    updateEmulator(emulator) {
      emulatorAPI.updateEmulator(emulator).then(_=>{
        message.success("修改成功");
        this.listEmulators();
      }).catch(resp=>{message.error(resp.message)});
    }
  }
}
</script>

<style scoped>
</style>