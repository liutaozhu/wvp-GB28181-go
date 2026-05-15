<template>
  <div id="operationsForSystemInfo" class="app-container system-info-page">
    <div class="system-info-content">
      <el-card v-for="(value, key) in systemInfoList" :key="key" class="info-card" shadow="hover">
        <div slot="header" class="card-header">
          <span>{{ key }}</span>
        </div>
        <el-descriptions :column="1" border label-class-name="info-label" content-class-name="info-content">
          <el-descriptions-item v-for="(childValue, childKey) in value" :key="childKey" :label="childKey">
            <span v-if="!childValue.startsWith('http')">{{ childValue }}</span>
            <a v-else target="_blank" :href="childValue">{{ childValue }}</a>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>
  </div>
</template>

<script>

export default {
  name: 'OperationsSystemInfo',
  data() {
    return {
      loading: false,
      systemInfoList: {}
    }
  },
  created() {
    this.initData()
  },
  methods: {
    initData: function() {
      this.loading = true
      this.$store.dispatch('server/info')
        .then(data => {
          this.systemInfoList = data
        })
        .catch((error) => {
          console.log(error)
        })
        .finally(() => {
          this.loading = false
        })
    }
  }
}
</script>

<style scoped>
.system-info-page {
  padding: 20px;
}

.system-info-content {
  max-width: 800px;
  margin: 20px auto;
}

.info-card {
  margin-bottom: 20px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
}
</style>

<style>
.info-label {
  width: 140px;
  text-align: right;
}

.info-content {
  text-align: left;
}
</style>
