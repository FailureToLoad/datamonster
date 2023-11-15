import { useLocation } from "react-router-dom";
import { SettlementContext } from "./context";
import { useContext, useEffect } from "react";
import Spinner from "@/components/spinner";

function Settlement() {
  const { settlement, setSettlement, loading, setLoading } =
    useContext(SettlementContext);
  const passedSettlement = useLocation().state.settlement;

  useEffect(() => {
    setSettlement(passedSettlement);
    setLoading(false);
  }, [setSettlement, setLoading]);
  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      {loading ? <Spinner /> : <h1>Settlement Id: {settlement?.id}</h1>}
    </div>
  );
}

export default Settlement;
