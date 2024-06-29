import { Link, Navigate, Outlet, useParams } from "react-router-dom";
import { navigationMenuTriggerStyle } from "@/components/ui/navigation-menu";
import { Settlement } from "@/types";
import { useAuth0 } from "@auth0/auth0-react";
import { Get } from "@/api";
import Spinner from "@/components/spinner";
import { useQuery } from "@tanstack/react-query";

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
  const { settlementId } = useParams();
  const { getAccessTokenSilently } = useAuth0();
  const getSettlement = async () => {
    try {
      let token = await getAccessTokenSilently();
      let response = await Get<Settlement | null>(
        "settlement/" + settlementId,
        token,
      );
      if (!response.data) {
        return null;
      }
      return response.data as Settlement;
    } catch (e) {
      console.log(e);
      return null;
    }
  };
  const { isPending, isError, data, error } = useQuery({
    queryKey: ["settlement"],
    queryFn: getSettlement,
  });

  if (isPending) {
    return <Spinner />;
  }

  if (isError) {
    return <span>Error: {error.message}</span>;
  }

  if (!data) {
    return <Navigate to="/" />;
  }

  let settlement = data as Settlement;
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
