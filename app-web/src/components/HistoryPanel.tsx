import { List, Typography } from "antd";
import type { HistoryEntry } from "@/types/calculator";

const { Text } = Typography;

interface HistoryPanelProps {
  history: HistoryEntry[];
  onClear: () => void;
}

export default function HistoryPanel({ history, onClear }: HistoryPanelProps) {
  if (history.length === 0) return null;

  return (
    <div className="history">
      <div className="history-header">
        <Text strong>History</Text>
        <Text
          type="secondary"
          className="history-clear"
          onClick={onClear}
        >
          Clear
        </Text>
      </div>
      <List
        size="small"
        dataSource={history}
        renderItem={(entry) => (
          <List.Item className="history-item">
            <Text type="secondary">{entry.expression} =</Text>
            <Text strong className="history-result">
              {entry.result}
            </Text>
          </List.Item>
        )}
      />
    </div>
  );
}
