import { ConfigProvider } from "antd";
import Calculator from "@/components/Calculator";

import "@/styles/index.css";

export default function App() {
  return (
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: "#018da2",
          borderRadius: 6,
          fontFamily:
            "'Archivo', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
        },
      }}
    >
      <div className="app-container">
        <Calculator />
      </div>
    </ConfigProvider>
  );
}
