import {
  ArchiveIcon,
  PersonIcon,
  type IconProps,
  HourglassMediumIcon,
} from "@phosphor-icons/react";
import { Link, Outlet, useLocation } from "react-router";


function LeftNav() {
  const { pathname } = useLocation();
  const timelineKey = "timeline";
  const populationKey = "population";
  const storageKey = "storage";
  const getProps = (active: boolean): IconProps => {
    const props: IconProps = {
      size: 32,
    };
    if (active) {
      props.weight = "fill";
      props.className = "text-primary";
    }
    return props;
  };

  return(
    <div className="h-screen absolute grid top-0 left-0 p-4">
        <ul className="menu bg-base-300 rounded-box mt-6">
          <li>
            <Link to={timelineKey} color="foreground" className="tooltip tooltip-right" data-tip="Timeline">
              <HourglassMediumIcon {...getProps(pathname.includes(timelineKey))} />
            </Link>
          </li>
          <li>
            <Link to={populationKey} color="foreground" className="tooltip tooltip-right" data-tip="Population">
              <PersonIcon {...getProps(pathname.includes(populationKey))} />
            </Link>
          </li>
          <li>
            <Link to={storageKey} color="foreground" className="tooltip tooltip-right" data-tip="Storage">
              <ArchiveIcon {...getProps(pathname.includes(storageKey))} />
            </Link>
          </li>
        </ul>
      </div>
  );
}

export default function SettlementPage() {
  return (
    <>
      <LeftNav />
      <div className="flex h-screen w-full flex-col justify-center overflow-auto">
        <div className="p-16 flex flex-1 justify-center">
          <Outlet />
        </div>
      </div>
    </>
  );
}
