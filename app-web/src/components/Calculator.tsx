import { Card, Spin } from "antd";
import Display from "@/components/Display";
import Keypad from "@/components/Keypad";
import HistoryPanel from "@/components/HistoryPanel";
import { useCalculator } from "@/hooks/useCalculator";

export default function Calculator() {
  const {
    display,
    expression,
    history,
    error,
    loading,
    pendingOp,
    appendDigit,
    setOperation,
    calculate,
    clear,
    clearAll,
    toggleSign,
    backspace,
  } = useCalculator();

  return (
    <Spin spinning={loading} size="small">
      <Card className="calculator" variant="borderless">
        <Display value={display} expression={expression} error={error} />
        <Keypad
          onDigit={appendDigit}
          onOperation={setOperation}
          onEquals={calculate}
          onClear={clear}
          onClearAll={clearAll}
          onToggleSign={toggleSign}
          onBackspace={backspace}
          pendingOp={pendingOp}
          loading={loading}
        />
        <HistoryPanel history={history} onClear={clearAll} />
      </Card>
    </Spin>
  );
}
