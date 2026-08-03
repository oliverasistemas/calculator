import { Button } from "antd";
import {
  DeleteOutlined,
  PlusOutlined,
  MinusOutlined,
  CloseOutlined,
  PercentageOutlined,
} from "@ant-design/icons";
import type { Operation } from "@/types/calculator";

interface KeypadProps {
  onDigit: (digit: string) => void;
  onOperation: (op: Operation) => void;
  onEquals: () => void;
  onClear: () => void;
  onClearAll: () => void;
  onToggleSign: () => void;
  onBackspace: () => void;
  pendingOp: Operation | null;
  loading: boolean;
}

const OP_LABELS: Record<string, React.ReactNode> = {
  add: <PlusOutlined />,
  subtract: <MinusOutlined />,
  multiply: <CloseOutlined />,
  divide: "\u00F7",
  power: "x\u207F",
  sqrt: "\u221A",
  percentage: <PercentageOutlined />,
};

export default function Keypad({
  onDigit,
  onOperation,
  onEquals,
  onClear,
  onClearAll,
  onToggleSign,
  onBackspace,
  pendingOp,
  loading,
}: KeypadProps) {
  const opButton = (op: Operation) => (
    <Button
      className={`key key-op ${pendingOp === op ? "key-active" : ""}`}
      onClick={() => onOperation(op)}
      disabled={loading}
    >
      {OP_LABELS[op]}
    </Button>
  );

  const digitButton = (d: string, className?: string) => (
    <Button
      className={`key ${className ?? ""}`}
      onClick={() => onDigit(d)}
      disabled={loading}
    >
      {d}
    </Button>
  );

  return (
    <div className="keypad">
      <Button className="key key-fn" onClick={onClearAll} disabled={loading}>
        AC
      </Button>
      <Button className="key key-fn" onClick={onClear} disabled={loading}>
        C
      </Button>
      <Button className="key key-fn" onClick={onBackspace} disabled={loading}>
        <DeleteOutlined />
      </Button>
      {opButton("divide")}

      {opButton("sqrt")}
      {opButton("power")}
      {opButton("percentage")}
      {opButton("multiply")}

      {digitButton("7")}
      {digitButton("8")}
      {digitButton("9")}
      {opButton("subtract")}

      {digitButton("4")}
      {digitButton("5")}
      {digitButton("6")}
      {opButton("add")}

      {digitButton("1")}
      {digitButton("2")}
      {digitButton("3")}
      <Button
        className="key key-equals"
        type="primary"
        onClick={onEquals}
        disabled={loading}
      >
        =
      </Button>

      <Button
        className="key key-fn"
        onClick={onToggleSign}
        disabled={loading}
      >
        +/-
      </Button>
      {digitButton("0")}
      {digitButton(".")}
      <div />
    </div>
  );
}
