import { BrowserRouter, Route, Routes, Outlet } from "react-router-dom";
import { AppShell } from "./components/layout/AppShell";
import { Inventory } from "./pages/Inventory";
import { QuoteBuilder } from "./pages/QuoteBuilder";
import OrderList from "./pages/orders/OrderList";
import OrderDetail from "./pages/orders/OrderDetail";

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
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
