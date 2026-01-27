import type { MessageApi, DialogApi, NotificationApi, LoadingBarApi } from 'naive-ui'

declare global {
  interface Window {
    $message: MessageApi
    $dialog: DialogApi
    $notification: NotificationApi
    $loadingBar: LoadingBarApi
  }

  const $message: MessageApi
  const $dialog: DialogApi
  const $notification: NotificationApi
  const $loadingBar: LoadingBarApi
}
