
<template>
  <div class="table-box">
    <div class="title">
      <h2>最简单的 CRUD Demo</h2>
    </div>
    <div class="query-box">
      <el-input class="query-input" v-model="queryInput" placeholder="请输入姓名搜索🔍" @change="handleQueryName" />
      <div class="btn-list">
        <el-button type="primary" @click="handleAdd">新增</el-button>
        <el-button type="primary" @click="handleDelList" v-if="multipleSelection.length > 0">删除多选</el-button>
      </div>
    </div>


    <el-table ref="multipleTableRef" :data="tableData" style="width: 100%" border
      @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="姓名" width="120" />
      <el-table-column prop="email" label="邮箱" width="120" />
      <el-table-column prop="phone" label="电话" width="120" />
      <el-table-column prop="state" label="状态" width="120" />
      <el-table-column prop="address" label="地址" width="150" />

      <el-table-column fixed="right" label="操作" width="120">
        <template #default="scope">
          <el-button link type="primary" size="small" @click="handleRowDel(scope.row.ID)" style="color: red">
            删除
          </el-button>
          <el-button link type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>

        </template>
      </el-table-column>
    </el-table>

    <el-pagination background layout="prev, pager, next" style="display: flex;justify-content: center;margin-top:10px; "
      :total="total" v-model:current-page="curPage" @current-change="handleChangePage" />


    <el-dialog v-model="dialogFormVisible" :title="dialogType === 'add' ? '新增' : '编辑'">


      <el-form :model="tableForm">
        <el-form-item label="姓名" :label-width="100">
          <el-input v-model="tableForm.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="邮箱" :label-width="100">
          <el-input v-model="tableForm.email" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="电话" :label-width="100">
          <el-input v-model="tableForm.phone" autocomplete="off" />
        </el-form-item>
        <el-form-item label="状态" :label-width="100">
          <el-input v-model="tableForm.state" autocomplete="off" />
        </el-form-item>
        <el-form-item label="地址" :label-width="100">
          <el-input v-model="tableForm.address" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="dialogConfirm">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>

  </div>
</template>




<style scoped>
.table-box {
  width: 800px;
  margin: 200px auto;
}

.title {
  text-align: center;
}

.query-box {
  display: flex;
  justify-content: space-between;
}

.query-input {
  width: 200px;
  margin-bottom: 20px;
}
</style>



<script setup >
import { ref } from "vue";
import request from "./utils/request.js"


// 数据
let queryInput = ref("")
let tableData = ref([
  {
    id: '1',
    name: 'Tom1',
    email: '123@163.com',
    phone: '18800115555',
    state: 'California',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    id: '2',
    name: 'Tom2',
    email: '123@163.com',
    phone: '18800115555',
    state: 'California',
    address: 'No. 189, Grove St, Los Angeles',
  },
  {
    id: '3',
    name: 'Tom3',
    email: '123@163.com',
    phone: '18800115555',
    state: 'California',
    address: 'No. 189, Grove St, Los Angeles',
  },
])

let multipleSelection = ref([])
let dialogFormVisible = ref(false)
let tableForm = ref({
  name: '郭心玥',
  email: "x@163.com",
  phone: "18800001111",
  state: "我才大二",
  address: "天水"
})
let dialogType = ref('add')


let tableDataCopy = Object.assign(tableData.value)

let total = ref(10)
let curPage = ref(1)



// 方法

const getTableData = async (cur = 1) => {
  // 第一种请求方式
  // let res = await request.get('list/?pageSize=10&pageNum=${cur}')
  // console.log(res)


  // 第二种请求方式
  let res = await request.get('/list/', {
    pageSize: 10,
    pageNum: cur,
  })
  // console.log(res)

  tableData.value = res.list
  total.value = res.total
  curPage.value = res.pageNum
}

getTableData(1)

// 请求分页
const handleChangePage = () => {
  getTableData(curPage.value)
}

// 搜索
const handleQueryName = async (val) => {
  // 相同
  // console.log(queryInput.value)
  console.log(val)

  // 精准查询
  // if (val.length > 0) {
  //   tableData.value = tableData.value.filter(item => (item.name).toLowerCase().match(val.toLowerCase()))
  // } else {
  //   tableData.value = tableDataCopy
  // }

  // 模糊查询
  if (val.length > 0) {
    tableData.value = await request.get('/list/' + val)
  } else {
    await getTableData(curPage.value)
  }
}

// 编辑
const handleEdit = (row) => {
  dialogFormVisible.value = true
  dialogType = 'edit'
  tableForm.value = { ...row }
  // console.log(row)
}


// 删除1条
const handleRowDel = async (id) => {
  // console.log(ID)
  // let index = tableData.value.findIndex(item => item.id === id)
  // tableData.value.splice(index, 1)

  // console.log(id)
  // console.log(index)

  await request.delete(`/delete/` + id)
  await getTableData(curPage.value)
}

// 删除多条
const handleDelList = async () => {
  multipleSelection.value.forEach(id => {
    handleRowDel(id)
  })
  // console.log(multipleSelection.value)
  multipleSelection.value = []
  // console.log(multipleSelection.value)
}

// 选中
const handleSelectionChange = (val) => {
  // multipleSelection.value = val
  // console.log(val)

  multipleSelection.value = []
  val.forEach(item => {
    multipleSelection.value.push(item.ID)
  })
  // console.log(multipleSelection.value)
}


// 新增
const handleAdd = () => {
  dialogType = 'add'
  dialogFormVisible.value = true
  tableForm.value = {}
  console.log(dialogFormVisible.value)
}



// 确认
const dialogConfirm = async () => {
  dialogFormVisible.value = false // 关闭弹窗


  // 判断是新增还是编辑
  if (dialogType === 'add') {
    // 1. 拿到数据
    // console.log(tableData.value);
    // 2. 添加到 table
    // 前端操作
    // tableData.value.push({
    //   id: (tableData.value.length + 1).toString(),
    //   ...tableForm.value,
    // })

    // 添加数据
    await request.post("/add", { ...tableForm.value })
    // 刷新数据
    await getTableData(curPage.value)

  } else {
    // 1. 获取到当前这条的索引
    // let index = tableData.value.findIndex(item => item.id === tableForm.value.id)
    // console.log(index);

    // 2. 替换当前索引值对应的数据
    // tableData.value[index] = tableForm.value

    // 1. 获取到当前这条的索引
    await request.put(`/update/` + tableForm.value.ID, { ...tableForm.value })

    // 2. 替换当前索引值对应的数据
    // 刷新数据
    await getTableData(curPage.value)

  }
}

</script>




