import { Link, Outlet, useLoaderData } from "react-router-dom";
import { navigationMenuTriggerStyle } from "@/components/ui/navigation-menu";
import { Settlement } from "@/api/settlement";

interface HeaderProps {
  settlement: Settlement;
}
function Header({ settlement }: HeaderProps) {
  return (
    <div id="header" className="sticky top-0 w-full flex-none">
      <div className="m-2 flex h-20 justify-center">
        <div className="flex w-1/3 flex-col">
          <div className="flex w-full">
            <div className="inline-flex grow justify-center">
              {settlement?.name}
            </div>
          </div>
          <div className="inline-flex justify-center">
            <Link to="timeline/" className={navigationMenuTriggerStyle()}>
              Timeline
            </Link>
            <Link to="population/" className={navigationMenuTriggerStyle()}>
              Population
            </Link>
            <Link to="storage/" className={navigationMenuTriggerStyle()}>
              Storage
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

function SettlementPage() {
  const settlement = useLoaderData() as Settlement;
  return (
    <div className="flex h-screen w-full flex-col">
      <Header settlement={settlement as Settlement} />
      <div className="p-16">
        <Outlet />
      </div>
    </div>
  );
}

export default SettlementPage;
