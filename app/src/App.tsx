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
import { DriverLayout } from "./pages/driver/DriverLayout";
import { RouteList } from "./pages/driver/RouteList";
import { StopList } from "./pages/driver/StopList";
import { DeliveryDetail } from "./pages/driver/DeliveryDetail";
import { DoorConfigurator } from "./pages/millwork/DoorConfigurator";
import { PartnerLayout } from "./components/layout/PartnerLayout";
import { PartnerDashboard } from "./pages/partner/Dashboard";
import { ProjectList } from "./pages/partner/ProjectList";
import { PartnerInvoiceList } from "./pages/partner/InvoiceList";
import { RFCDashboard } from "./pages/governance/RFCDashboard";
import { NewRFC } from "./pages/governance/NewRFC";
import { RFCDetail } from "./pages/governance/RFCDetail";

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
          <Route path="millwork/configure" element={<DoorConfigurator />} />
          <Route path="governance">
            <Route index element={<RFCDashboard />} />
            <Route path="new" element={<NewRFC />} />
            <Route path=":id" element={<RFCDetail />} />
          </Route>
        </Route>

        {/* Mobile Driver App */}
        <Route path="/driver" element={<DriverLayout />}>
          <Route index element={<RouteList />} />
          <Route path="routes/:id" element={<StopList />} />
          <Route path="deliveries/:id" element={<DeliveryDetail />} />
        </Route>

        {/* Partner Portal */}
        <Route path="/partner" element={<PartnerLayout />}>
          <Route index element={<PartnerDashboard />} />
          <Route path="projects" element={<ProjectList />} />
          <Route path="invoices" element={<PartnerInvoiceList />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
