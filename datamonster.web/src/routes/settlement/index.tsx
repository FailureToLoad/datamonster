import { Link, Navigate, Outlet, useLocation } from "react-router-dom";
import { Separator } from "@/components/ui/separator";
import {
  NavigationMenu,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import { Settlement } from "@/api/api";

interface HeaderProps {
  settlement: Settlement;
}
function Header({ settlement }: HeaderProps) {
  return (
    <div id="header" className="sticky top-0 w-full flex-none">
      <div className="m-2 flex h-20">
        <div className="flex items-center">
          <div className="mx-2 border-2 border-solid border-black px-4 py-2 text-lg font-semibold">
            {settlement?.limit}
          </div>
          <div>Survival Limit</div>
        </div>
        <Separator
          className="m-4 w-[1px] shrink-0 bg-black"
          orientation="vertical"
        />
        <div className="flex w-2/6 flex-col items-start justify-center">
          <div className="flex w-full">
            <h1 className="w-auto flex-none font-medium">Settlement Name</h1>
            <div className="inline-flex grow justify-center">
              {settlement?.name}
            </div>
          </div>
          <Separator className="my-1 h-[2px] w-full shrink-0 bg-black " />
          <div>
            <NavigationMenu>
              <NavigationMenuList>
                <Link to="timeline/" className={navigationMenuTriggerStyle()}>
                  Timeline
                </Link>
                <Link to="population/" className={navigationMenuTriggerStyle()}>
                  Population
                </Link>
                <Link to="storage/" className={navigationMenuTriggerStyle()}>
                  Storage
                </Link>
              </NavigationMenuList>
            </NavigationMenu>
          </div>
        </div>
      </div>
    </div>
  );
}

function Settlement() {
  const settlementJson = localStorage.getItem("settlement");
  const settlement = settlementJson
    ? (JSON.parse(settlementJson) as Settlement)
    : null;
  if (!settlement) {
    return <Navigate to="/select" />;
  }
  return (
    <div className="flex h-screen w-full flex-col">
      <Header settlement={settlement as Settlement} />
      <Outlet />
    </div>
  );
}

export default Settlement;
