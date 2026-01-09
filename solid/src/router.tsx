
import { Router, Route, Navigate } from '@solidjs/router';
import ProtectedLayout, { checkAuth } from '~/components/ProtectedLayout';
import SelectSettlement, { getSettlements } from '~/pages/selection/page';


export function AppRouter() {
  return (
    <Router>
      <Route path="/" component={ProtectedLayout} preload={() => checkAuth()}>
        <Route path="/" component={() => <Navigate href="/settlements" />} />
        <Route path="/settlements" component={SelectSettlement} preload={() => getSettlements()} />
        <Route path="/settlements/:settlementId">
          <Route path="/" component={() => <Navigate href="timeline" />} />
          <Route path="/population" component={() => <div>Population</div>} />
          <Route path="/storage" component={() => <div>Storage</div>} />
          <Route path="/timeline" component={() => <div>Timeline</div>} />
        </Route>
      </Route>
    </Router>
  );
}
