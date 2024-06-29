import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Get } from "@/api";
import { Settlement } from "@/types";
import { CreateSettlementDialogue } from "./createSettlementDialogue";
import { Button } from "@/components/ui/button";
import { Play } from "lucide-react";
import { Link } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import Spinner from "@/components/spinner";
import { useAuth0 } from "@auth0/auth0-react";

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

type SelectorProps = {
  testId?: string;
};

function SettlementSelector({ testId }: SelectorProps) {
  const { getAccessTokenSilently } = useAuth0();
  const getSettlements = async () => {
    try {
      const token = await getAccessTokenSilently();
      const response = await Get<Array<Settlement> | null>("settlement", token);
      if (!response.data) {
        return null;
      }
      return response.data;
    } catch (e) {
      console.log(e);
      return null;
    }
  };

  const { isPending, isError, data, error } = useQuery({
    queryKey: ["settlements"],
    queryFn: getSettlements,
  });

  if (isPending) {
    return <Spinner />;
  }

  if (isError) {
    throw new Error(error.message);
  }
  console.log(data);
  let settlements = data as Array<Settlement>;
  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden">
      <ul className="w-1/4 space-y-4 " data-testid={testId}>
        {settlements &&
          settlements.map((settlement) => (
            <li key={settlement.id}>
              <SettlementCard settlement={settlement} />
            </li>
          ))}
        <li key={-1}>
          <CreateSettlementDialogue />
        </li>
      </ul>
    </div>
  );
}

export default SettlementSelector;
