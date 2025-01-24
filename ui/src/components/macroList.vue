<template>
  <a-row>
    <a-select :options="emulatorTypeOptions" v-model:value="selectedEmulatorType" @change="listMacros"></a-select>
    <a-button type="primary" @click="openMacroModal">新建</a-button>
  </a-row>
  <a-table :columns="columns" :data-source="macros" :pagination="false">
    <template #bodyCell="{ column, text, record }">
      <template v-if="column.dataIndex === 'action'">
        <a-button type="primary" v-if="showAction" @click="applyMacro(record)">使用</a-button>
        <a-button @click="removeMacro(record['macroId'])">删除</a-button>
      </template>
      <template v-if="column.dataIndex === 'keyCodes'">
        {{record.keyCodes}}
      </template>
    </template>
  </a-table>
  <a-modal :open="createMacroModalOpen" title="创建按键宏" @cancel="closeMacroModal" @ok="createMacro">
    <a-form>
      <a-form-item label="模拟器类型">
        {{selectedEmulatorType}}
      </a-form-item>
      <a-form-item label="名称">
        <a-input v-model:value="newMacroName"></a-input>
      </a-form-item>
      <a-form-item label="按键组合(点击删除)">
        <a-row>
          <a-button v-for="btn in selectedButtons" @click="removeButton(btn)">
            {{btn}}
          </a-button>
        </a-row>
      </a-form-item>
      <a-form-item label="可选按键(点击添加)">
        <a-row>
          <a-button v-for="btn in emulatorButtons" @click="addButton(btn)">
            {{btn.label}}
          </a-button>
        </a-row>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import {Row, Col, Card, Modal, Select as ASelect, Checkbox as ACheckbox, message} from "ant-design-vue";
import { Button, Table, Form, Input} from "ant-design-vue";
import api from "../api/request.js";
import emulatorAPI from "../api/emulator.js";
import configs from "../api/const.js";

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

      createMacroModalOpen: false,
      emulatorButtons: [],
      selectedButtons: [],
      newMacroName: "",
    }
  },
  created() {
    addEventListener("emulator_restart", ev=>{
      this.selectedEmulatorType = ev.detail["emulatorType"];
    });
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
    openMacroModal() {
      const buttons = configs.defaultButtonLayouts[this.selectedEmulatorType];
      if (buttons) {
        this.emulatorButtons = buttons;
        this.createMacroModalOpen = true;
      }else {
        message.info("该模拟器没有按键");
      }
    },
    closeMacroModal() {
      this.createMacroModalOpen = false;
      this.selectedButtons = [];
      this.emulatorButtons = [];
    },
    addButton(btn) {
      btn["codes"].forEach(code=>{
        if (!this.selectedButtons.includes(code)) {
          this.selectedButtons.push(code);
        }
      });
    },
    removeButton(btn) {
      this.selectedButtons = this.selectedButtons.filter(button => button !== btn);
    },
    createMacro() {
      if (this.newMacroName.length === 0) {
        message.warn("请输入宏名称");
        return;
      }
      if (this.newMacroName.length >= 16) {
        message.warn("名称太长");
        return;
      }
      if (this.selectedButtons.length === 0) {
        message.warn("请选择按键组合");
        return;
      }
      api.post("/macros", {
        macroName: this.newMacroName,
        emulatorType: this.selectedEmulatorType,
        keyCodes: this.selectedButtons,
      }).then(resp=>{
        message.success(resp["message"]);
        this.listMacros();
        this.closeMacroModal();
      }).catch(err=>{
        message.error(err["message"]);
      });
    },
    removeMacro(id) {
      api.delete("/macros/" + id).then(resp=>{
        message.success(resp["message"]);
        this.listMacros();
      }).catch(err=>{
        message.error(err["message"]);
      })
    },
  }
}
</script>

<style scoped>
</style>