import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";

import api, { Settlement } from "@/api/api";
import { useCallback, useContext, useEffect, useState } from "react";
import { AuthContext } from "@/auth/auth-context";
import Spinner from "@/components/spinner";
import { CreateSettlementDialogue } from "./createSettlementDialogue";

interface SettlementListProps {
  token: string;
  update: (s: Settlement) => void;
  settlements: Array<Settlement>;
}

function SettlementList({ token, update, settlements }: SettlementListProps) {
  const dialogueListItem = (
    <li key={-1}>
      <CreateSettlementDialogue token={token} update={update} />
    </li>
  );
  if (settlements.length < 1) {
    return dialogueListItem;
  }

  const cards = settlements.map((settlement) => (
    <li key={settlement.id}>
      <Card>
        <CardHeader>
          <CardTitle>{settlement.name}</CardTitle>
          <CardDescription>Lantern Year: {settlement.year}</CardDescription>
        </CardHeader>
      </Card>
    </li>
  ));
  return [dialogueListItem, ...cards];
}

function SettlementSelector() {
  const [settlements, setSettlements] = useState<Array<Settlement>>(
    Array<Settlement>(),
  );
  const [isLoading, setIsLoading] = useState(true);
  const [token, setToken] = useState("");
  const { currentUser } = useContext(AuthContext);
  const updateSettlementList = useCallback(
    (s: Settlement) => {
      setSettlements([s, ...settlements]);
    },
    [settlements],
  );

  useEffect(() => {
    if (!currentUser) {
      return;
    }
    if (!token) {
      currentUser.getIdToken().then((idToken) => {
        setToken(idToken);
      });
    } else {
      api.getSettlementsForUser(token).then((val) => {
        setSettlements(val);
        setIsLoading(false);
      });
    }
  }, [currentUser, token]);

  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      {isLoading ? (
        <Spinner />
      ) : (
        <ul className="w-1/4 space-y-4 ">
          <SettlementList
            token={token}
            update={updateSettlementList}
            settlements={settlements}
          />
        </ul>
      )}
    </div>
  );
}

export default SettlementSelector;
