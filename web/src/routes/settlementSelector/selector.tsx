import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Settlement } from "@/api/settlement";
import { useCallback, useState } from "react";
import { CreateSettlementDialogue } from "./createSettlementDialogue";
import { Button } from "@/components/ui/button";
import { Play } from "lucide-react";
import { Link, useLoaderData } from "react-router-dom";

function SettlementCard({ settlement }: { settlement: Settlement }) {
  const link = "/" + settlement.id;
  return (
    <Card>
      <CardHeader>
        <CardTitle>{settlement.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-row justify-between">
          <div>Lantern Year: {settlement.year}</div>
          <div>
            <Link to={link}>
              <Button variant="ghost" size="icon">
                <Play className="h-6 w-6" />
              </Button>
            </Link>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

interface SettlementListProps {
  update: (s: Settlement) => void;
  settlements: Array<Settlement>;
}

function SettlementList({ update, settlements }: SettlementListProps) {
  const dialogueListItem = (
    <li key={-1}>
      <CreateSettlementDialogue update={update} />
    </li>
  );
  if (settlements === null) {
    return dialogueListItem;
  }

  const cards = settlements.map((settlement) => (
    <li key={settlement.id}>
      <SettlementCard settlement={settlement} />
    </li>
  ));
  return [dialogueListItem, ...cards];
}

type SelectorProps = {
  testId?: string;
};

function SettlementSelector({ testId }: SelectorProps) {
  const loadedSettlements = useLoaderData() as Array<Settlement>;
  const [settlements, setSettlements] =
    useState<Array<Settlement>>(loadedSettlements);
  const updateSettlementList = useCallback(
    (s: Settlement) => {
      if (!settlements || settlements.length === 0) {
        setSettlements([s]);
        return;
      }
      setSettlements([s, ...settlements]);
    },
    [settlements],
  );

  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      <ul className="w-1/4 space-y-4 " data-testid={testId}>
        <SettlementList
          update={updateSettlementList}
          settlements={settlements}
        />
      </ul>
    </div>
  );
}

export default SettlementSelector;
