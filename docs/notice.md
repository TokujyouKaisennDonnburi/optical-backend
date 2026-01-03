# Noticeについて

## creatorについて
NoticeのCreateは他と違い、API化せず、ほか関数から呼び出す形になります。

## CreateNotice の使い方

### 必要
略/internal/notice/creator"のインポート

### 関数の説明
- 必須パラメータ

| パラメータ | 型 | 説明 |
| - | - | - |
| ctx        | context.Context | コンテキスト     |
| userID     | uuid.UUID       | 通知先ユーザーID |
| title      | string          | 通知タイトル     |
| content    | string          | 通知本文         |

- オプションパラメータ

| パラメータ | 型 | 説明 |
| - | - | - |
| creator.WithEventID(id)    | uuid.UUID | イベントID   |
| creator.WithCalendarID(id) | uuid.UUID | カレンダーID |


