<template>
  <a-row>
    <a-select :options="emulatorTypeOptions" v-model:value="selectedEmulatorType" @change="listMacros"></a-select>
    <a-button type="primary">新建</a-button>
  </a-row>
  <a-table :columns="columns" :data-source="macros" :pagination="false">
    <template #bodyCell="{ column, text, record }">
      <template v-if="column.dataIndex === 'action'">
        <a-button type="primary" :hidden="!showAction" @click="applyMacro(record)">使用</a-button>
      </template>
      <template v-if="column.dataIndex === 'keyCodes'">
        {{record.keyCodes}}
      </template>
    </template>
  </a-table>
</template>

<script>
import {Row, Col, Card, Modal, Select as ASelect, Checkbox as ACheckbox} from "ant-design-vue";
import { Button, Table, Form, Input} from "ant-design-vue";
import api from "../api/request.js";
import emulatorAPI from "../api/emulator.js";

export default {
  components: {
    ACheckbox,
    ASelect,
    ARow: Row,
    ACol: Col,
    AButton: Button,
    ACard: Card,
    ATable: Table,
    AModal: Modal,
    AFormItem: Form.Item,
    AForm: Form,
    AInput: Input,
  },
  props: {
    showAction: Boolean,
    roomId: String,
  },
  data() {
    return {
      macros: [],
      columns: [
        {title: "宏名称", dataIndex: "macroName"},
        {title: "按钮", dataIndex: "keyCodes"},
        {title: "操作", dataIndex: "action"},
      ],
      emulatorTypes: [],
      selectedEmulatorType: null,
      emulatorTypeOptions: [],
    }
  },
  created() {
    emulatorAPI.listEmulatorTypes().then(resp=>{
      this.emulatorTypes = resp.data;
      this.selectedEmulatorType = this.emulatorTypes[0];
      let options = [];
      this.emulatorTypes.forEach(emulatorType => {
        options.push({label: emulatorType, value: emulatorType});
      })
      this.emulatorTypeOptions = options;
      this.listMacros();
    });
  },
  unmounted() {
  },
  methods: {
    listMacros() {
      api.get("/macros", {emulatorType: this.selectedEmulatorType}).then(resp=>{
        this.macros = resp.data;
      });
    },
    applyMacro(macro) {
      api.post("/macros/apply", {
        roomId: this.roomId,
        macroId: macro["macroId"],
      });
    },
  }
}
</script>

<style scoped>
</style>