<template>
  <n-modal
    v-model:show="showModal"
    preset="card"
    title="新建本地歌单"
    class="create-playlist-modal"
    :style="{ width: '400px' }"
    :bordered="false"
  >
    <n-form
      ref="formRef"
      :model="formValue"
      :rules="rules"
      label-placement="top"
    >
      <n-form-item label="歌单名称" path="title" required>
        <n-input v-model:value="formValue.title" placeholder="请输入歌单名称" />
      </n-form-item>
      <n-form-item label="歌单描述" path="description">
        <n-input
          v-model:value="formValue.description"
          type="textarea"
          placeholder="请输入歌单描述 (选填)"
        />
      </n-form-item>
    </n-form>
    <template #footer>
      <n-button type="primary" block @click="handleCreate" :loading="loading">
        新建
      </n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import { useMessage, NModal, NForm, NFormItem, NInput, NButton } from 'naive-ui';
import { createPrivatePlaylist } from '@/api/playlist';

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:show', 'created']);

const message = useMessage();
const showModal = ref(false);
const loading = ref(false);
const formRef = ref(null);

const formValue = reactive({
  title: '',
  description: ''
});

const rules = {
  title: {
    required: true,
    message: '请输入歌单名称',
    trigger: ['input', 'blur']
  }
};

watch(() => props.show, (val) => {
  showModal.value = val;
  if (val) {
    // Reset form when opening
    formValue.title = '';
    formValue.description = '';
  }
});

watch(showModal, (val) => {
  emit('update:show', val);
});

const handleCreate = (e: MouseEvent) => {
  e.preventDefault();
  (formRef.value as any)?.validate((errors: any) => {
    if (!errors) {
      loading.value = true;
      createPrivatePlaylist(formValue)
        .then((res: any) => {
          if (res.code === 200 || res.code === 1000) { // Assuming 1000 is success based on previous context
            message.success('创建成功');
            showModal.value = false;
            emit('created', res.data); // Pass back the new playlist
          } else {
            message.error(res.message || '创建失败');
          }
        })
        .catch((err: any) => {
            message.error(err.message || '创建失败');
        })
        .finally(() => {
          loading.value = false;
        });
    }
  });
};
</script>

<style scoped>
.create-playlist-modal {
    /* Custom styles if needed */
}
</style>
