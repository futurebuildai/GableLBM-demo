import { BrowserRouter, Route, Routes, Outlet } from "react-router-dom";
import { AppShell } from "./components/layout/AppShell";
import { Inventory } from "./pages/Inventory";
import { QuoteBuilder } from "./pages/QuoteBuilder";
import OrderList from "./pages/orders/OrderList";
import OrderDetail from "./pages/orders/OrderDetail";
import InvoiceList from "./pages/invoices/InvoiceList";
import InvoiceDetail from "./pages/invoices/InvoiceDetail";
import DailyTill from "./pages/DailyTill";
import { DispatchBoard } from "./pages/DispatchBoard";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<AppShell><Outlet /></AppShell>}>
          <Route index element={<div className="p-8 text-white">Dashboard Placeholder</div>} />
          <Route path="inventory" element={<Inventory />} />
          <Route path="quotes/new" element={<QuoteBuilder />} />
          <Route path="orders" element={<OrderList />} />
          <Route path="orders/:id" element={<OrderDetail />} />
          <Route path="invoices" element={<InvoiceList />} />
          <Route path="invoices/:id" element={<InvoiceDetail />} />
          <Route path="reports/daily-till" element={<DailyTill />} />
          <Route path="dispatch" element={<DispatchBoard />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
