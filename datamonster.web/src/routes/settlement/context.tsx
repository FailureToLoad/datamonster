import { Settlement } from "@/api/api";
import { ReactNode, createContext, useState } from "react";

interface Props {
  children?: ReactNode;
}

export const SettlementContext = createContext({
  settlement: {} as Settlement | null,
  setSettlement: (_settlement: Settlement) => {},
  loading: true,
  setLoading: (_loading: boolean) => {},
});

export const SettlementProvider = ({ children }: Props) => {
  const [settlement, setSettlement] = useState<Settlement | null>(null);
  const [loading, setLoading] = useState(true);
  const value = {
    settlement,
    setSettlement,
    loading,
    setLoading,
  };
  return (
    <SettlementContext.Provider value={value}>
      {children}
    </SettlementContext.Provider>
  );
};
