user.go

    役割: エンティティと型定義
    内容:
        User構造体
            基本情報 + IconURL（アイコンURL）
        リクエスト型
            CreateUserRequest - ユーザー登録用
            UpdateUserRequest - ユーザー更新用
            LoginRequest - ログイン用
            CreateProfileRequest - プロフィール作成用
            UpdateProfileRequest - プロフィール更新用
        レスポンス型
            UserResponse - 詳細ユーザー情報
            LoginResponse - ログイン結果
            ProfileResponse - プロフィール表示用
        変換メソッド
            ToResponse() - UserをUserResponseに変換
            ToProfileResponse() - UserをProfileResponseに変換

handler

    役割: HTTPエンドポイント
    主要メソッド:
        認証系
            Register() - ユーザー登録（POST /register）
            Login() - ログイン（POST /login）
            Logout() - ログアウト（POST /logout）
        プロフィール系
            CreateProfile() - プロフィール作成（POST /profiles）
            GetProfile() - プロフィール取得（GET /profiles/:profileId）
            UpdateProfile() - プロフィール編集（PATCH /profiles）

repository

    役割: データベース操作
    主要メソッド:
        Create() - ユーザー作成（アイコンURL含む）
        GetByID(), GetByEmail() - ユーザー取得
        Update() - ユーザー更新（アイコンURL含む）
        Delete() - 論理削除
        SearchByNameOrEmail() - ユーザー検索

service/command

    役割: 書き込みロジック
    主要機能:
        認証処理
            ユーザー登録（メール重複チェック、パスワードハッシュ化）
            ログイン認証（パスワード検証）
            ユーザー更新・削除
        プロフィール処理
            プロフィール作成
            プロフィール更新

service/query

    役割: 読み込みロジック
    主要機能:
        データ取得
            プロフィール取得
            ユーザー検索（ID、メールアドレス）
            ユーザー一覧取得
