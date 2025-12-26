package openrouter

import (
	"context"
	"errors"
)

const (
	reason_error          = "error"
	reason_tool_calls     = "tool_calls"
	reason_stop           = "stop"
	reason_length         = "length"
	reason_content_filter = "content_filter"
)

var (
	ErrFinishWithError         = errors.New("finished reason: error")
	ErrFinishWithLength        = errors.New("finished reason: length")
	ErrFinishWithContentFilter = errors.New("finished reason: content-filter")
)

func (r OpenRouter) ChainFetch(ctx context.Context, messages []Message) (*OpenRouterResponse, error) {
	resp, err := r.Fetch(ctx, messages)
	if err != nil {
		return nil, err
	}
	result := resp.Choices[0]
	// 終了理由によって分岐
	switch result.FinishReason {
	// 完了
	case reason_stop:
		return resp, nil
		// ツール呼び出し
	case reason_tool_calls:
		// 次に送信するメッセージに今回の結果を含める
		messages = append(messages, result.Message)
		// 各ツール呼び出しに対応
		for _, call := range result.Message.ToolCalls {
			// 名前からツールを取得
			tool, err := r.getToolByName(call.Function.Name)
			if err != nil {
				return nil, err
			}
			// ツール実行
			toolOutput, err := tool.Call(ctx, call.Function.Arguments)
			if err != nil {
				return nil, err
			}
			// 実行結果をメッセージに追加
			messages = append(messages, Message{
				Role:       "tool",
				ToolCallId: call.Id,
				Content:    toolOutput,
			})
		}
		return r.ChainFetch(ctx, messages)
		// エラー
	case reason_error:
		return nil, ErrFinishWithError
	case reason_length:
		return nil, ErrFinishWithLength
	case reason_content_filter:
		return nil, ErrFinishWithContentFilter
	default:
		return resp, nil
	}
}

func (r OpenRouter) ChainStream(
	ctx context.Context,
	messages []Message,
	streamingFn func(context.Context, []byte) error,
) (*OpenRouterResponse, error) {
	resp, err := r.Stream(ctx, messages, streamingFn)
	if err != nil {
		return nil, err
	}
	result := resp.Choices[0]
	// 終了理由によって分岐
	switch result.FinishReason {
	// 完了
	case reason_stop:
		return resp, nil
		// ツール呼び出し
	case reason_tool_calls:
		// 次に送信するメッセージに今回の結果を含める
		messages = append(messages, result.Message)
		// 各ツール呼び出しに対応
		for _, call := range result.Message.ToolCalls {
			// 名前からツールを取得
			tool, err := r.getToolByName(call.Function.Name)
			if err != nil {
				return nil, err
			}
			// ツール実行
			toolOutput, err := tool.Call(ctx, call.Function.Arguments)
			if err != nil {
				return nil, err
			}
			// 実行結果をメッセージに追加
			messages = append(messages, Message{
				Role:       "tool",
				ToolCallId: call.Id,
				Content:    toolOutput,
			})
		}
		streamingFn(ctx, []byte("{\"status\": \"generating\"}"))
		return r.ChainStream(ctx, messages, streamingFn)
		// エラー
	case reason_error:
		return nil, ErrFinishWithError
	case reason_length:
		return nil, ErrFinishWithLength
	case reason_content_filter:
		return nil, ErrFinishWithContentFilter
	}
	return resp, nil
}
