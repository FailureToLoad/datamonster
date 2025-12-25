import {
  ArchiveIcon,
  PersonIcon,
  type IconProps,
  HourglassMediumIcon,
  TentIcon,
  SignOutIcon,
} from "@phosphor-icons/react";
import { Link, Outlet, useLocation, useNavigate } from "react-router";


function LeftNav() {
  const { pathname } = useLocation();
  const navigate = useNavigate();
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

  const handleLogout = async () => {
    await fetch("/api/auth/logout", { credentials: "include" });
    navigate("/");
  };

  return (
    <div className="left-0 top-0 grid h-screen w-fit p-4">
      <ul className="menu mt-6 justify-between rounded-box bg-base-300">
        <div>
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
        </div>
        <div>
          <li>
            <Link to="/settlements" className="tooltip tooltip-right" data-tip="Settlements">
              <TentIcon size={32} />
            </Link>
          </li>
          <li>
            <button onClick={handleLogout} className="tooltip tooltip-right" data-tip="Logout">
              <SignOutIcon size={32} />
            </button>
          </li>
        </div>
      </ul>
    </div>
  );
}

export default function SettlementPage() {
  return (
    <div className="flex h-screen">
      <LeftNav />
      <div className="flex-1 overflow-auto p-16">
        <Outlet />
      </div>
    </div>
  );
}
