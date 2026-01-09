
import { Router, Route, Navigate } from '@solidjs/router';
import { lazy } from 'solid-js';
import ProtectedLayout, { checkAuth } from '~/components/ProtectedLayout';
import SelectSettlement, { getSettlements } from '~/pages/selection/page';
import { SettlementPage } from '~/pages/settlement/page';

const population = lazy(() => import("~/pages/settlement/population/tab"));
const timeline = lazy(() => import("~/pages/settlement/timeline/tab"));
const storage = lazy(() => import("~/pages/settlement/storage/tab"));

export function AppRouter() {
  return (
    <Router>
      <Route path="/" component={ProtectedLayout} preload={() => checkAuth()}>
        <Route path="/" component={() => <Navigate href="/settlements" />} />
        <Route path="/settlements" component={SelectSettlement} preload={() => getSettlements()} />
        <Route path="/settlements/:settlementId" component={SettlementPage}>
          <Route path="/" component={() => <Navigate href="timeline" />} />
          <Route path="/population" component={population} />
          <Route path="/storage" component={storage} />
          <Route path="/timeline" component={timeline} />
        </Route>
      </Route>
    </Router>
  );
}
