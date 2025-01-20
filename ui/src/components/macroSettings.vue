<template>
  <a-card :bordered="false" class="center-card">
    <macro-list :show-action="true"/>
  </a-card>
</template>

<script>
import {Row, Col, List, Card, Modal, UploadDragger} from "ant-design-vue";
import { Button, Table, Form, Select, Input, Pagination} from "ant-design-vue";
import MacroList from "./macroList.vue";

export default {
  components: {
    ARow: Row,
    ACol: Col,
    AButton: Button,
    AList: List,
    AListItem: List.Item,
    ACard: Card,
    ATable: Table,
    AModal: Modal,
    AUploadDragger: UploadDragger,
    AFormItem: Form.Item,
    ASelect: Select,
    AForm: Form,
    AInput: Input,
    APagination: Pagination,
    MacroList: MacroList,
  },
  data() {
    return {

    }
  },
  created() {

  },
  unmounted() {
  },
  methods: {
  }
}
</script>

<style scoped>
</style>