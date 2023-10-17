import {
  NavigationMenu,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import { Separator } from "@radix-ui/react-separator";
import api, { Settlement } from "../api/api";
import { useEffect, useState } from "react";

function Header() {
  const [settlement, setSettlement] = useState<Settlement>();
  const fetchSettlement = async () => {
    const data = await api.getSettlement();
    setSettlement(data);
  };
  useEffect(() => {
    fetchSettlement();
  }, []);
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
                <NavigationMenuLink
                  className={navigationMenuTriggerStyle()}
                  href="https://www.google.com"
                >
                  Timeline
                </NavigationMenuLink>
                <NavigationMenuLink
                  className={navigationMenuTriggerStyle()}
                  href="https://www.google.com"
                >
                  Population
                </NavigationMenuLink>
                <NavigationMenuLink
                  className={navigationMenuTriggerStyle()}
                  href="https://www.google.com"
                >
                  Storage
                </NavigationMenuLink>
              </NavigationMenuList>
            </NavigationMenu>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Header;
