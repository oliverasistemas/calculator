import { Alert } from "antd";

interface DisplayProps {
  value: string;
  expression: string;
  error: string | null;
}

export default function Display({ value, expression, error }: DisplayProps) {
  return (
    <div className="display">
      <div className="display-expression">{expression || "\u00A0"}</div>
      <div className="display-value">{value}</div>
      {error && (
        <Alert
          message={error}
          type="error"
          showIcon
          closable
          className="display-error"
        />
      )}
    </div>
  );
}
