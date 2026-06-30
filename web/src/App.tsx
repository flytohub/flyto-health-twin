import { useEffect, useMemo, useState } from "react";

type PublicRecord = {
  date: string;
  source_id?: string;
  sleep_score: number;
  sleep_hours: number;
  hrv: number;
  resting_heart_rate: number;
  steps: number;
  exercise_minutes: number;
  stress_score: number;
  fatigue_score: number;
  caffeine_servings: number;
  water_liters: number;
};

type PublicPrediction = {
  target_date: string;
  model_id: string;
  model_version: string;
  input_start_date: string;
  input_end_date: string;
  feature_set: string[];
  assumptions: string[];
  predicted_hrv: number;
  predicted_resting_heart_rate: number;
  predicted_fatigue_score: number;
  predicted_sleep_score: number;
  recovery_state: string;
  hints: string[];
  missing_variables: string[];
};

type PublicEvaluation = {
  target_date: string;
  prediction: PublicPrediction;
  actual: PublicRecord;
  hrv_error: number;
  rhr_error: number;
  fatigue_error: number;
  sleep_error: number;
};

type ErrorSummary = {
  hrv_mean_absolute_error: number;
  rhr_mean_absolute_error: number;
  fatigue_mean_absolute_error: number;
  sleep_mean_absolute_error: number;
};

type PublicExport = {
  project: string;
  generated_at: string;
  records: PublicRecord[];
  evaluations: PublicEvaluation[];
  error_summary: ErrorSummary;
};

type SeriesPoint = {
  label: string;
  actual: number;
  predicted: number;
};

const dataUrl = `${import.meta.env.BASE_URL}public-data.json`;
const logoUrl = `${import.meta.env.BASE_URL}brand/flyto-logo.png`;

export function App() {
  const [data, setData] = useState<PublicExport | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    fetch(dataUrl)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`Data request failed with ${response.status}`);
        }
        return response.json() as Promise<PublicExport>;
      })
      .then((payload) => {
        if (!cancelled) {
          setData(payload);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : "Unable to load public data");
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  if (error) {
    return <StatusScreen title="Public data unavailable" detail={error} />;
  }

  if (!data) {
    return <StatusScreen title="Loading public dashboard" detail="Fetching redacted daily aggregates." />;
  }

  return <Dashboard data={data} />;
}

function Dashboard({ data }: { data: PublicExport }) {
  const latestEvaluation = data.evaluations.at(-1);
  const latestRecord = data.records.at(-1);
  const model = latestEvaluation?.prediction;

  const hrvSeries = useMemo(
    () =>
      data.evaluations.map((item) => ({
        label: item.target_date,
        actual: item.actual.hrv,
        predicted: item.prediction.predicted_hrv,
      })),
    [data.evaluations],
  );

  const sleepSeries = useMemo(
    () =>
      data.evaluations.map((item) => ({
        label: item.target_date,
        actual: item.actual.sleep_score,
        predicted: item.prediction.predicted_sleep_score,
      })),
    [data.evaluations],
  );

  const fatigueSeries = useMemo(
    () =>
      data.evaluations.map((item) => ({
        label: item.target_date,
        actual: item.actual.fatigue_score,
        predicted: item.prediction.predicted_fatigue_score,
      })),
    [data.evaluations],
  );

  return (
    <main className="app-shell">
      <section className="top-band">
        <div>
          <div className="brand-row">
            <img src={logoUrl} alt="Flyto2 logo" />
            <p className="eyebrow">Public research dashboard</p>
          </div>
          <h1>Flyto2</h1>
          <p className="lede">
            Privacy-preserving daily aggregates, transparent next-day predictions, and error analysis
            for a personal health digital twin prototype.
          </p>
        </div>
        <div className="status-panel" aria-label="Dataset status">
          <span className="status-dot" />
          <span>{data.records.length} synthetic daily records</span>
        </div>
      </section>

      <section className="metric-grid" aria-label="Summary metrics">
        <MetricCard label="Latest recovery" value={formatState(model?.recovery_state)} tone={model?.recovery_state} />
        <MetricCard label="HRV MAE" value={formatNumber(data.error_summary.hrv_mean_absolute_error)} unit="ms" />
        <MetricCard label="RHR MAE" value={formatNumber(data.error_summary.rhr_mean_absolute_error)} unit="bpm" />
        <MetricCard label="Sleep MAE" value={formatNumber(data.error_summary.sleep_mean_absolute_error)} unit="score" />
      </section>

      <section className="dashboard-grid">
        <div className="chart-panel wide">
          <PanelHeader title="Prediction error loop" meta={`${data.evaluations.length} evaluated days`} />
          <LineChart points={hrvSeries} actualLabel="Actual HRV" predictedLabel="Predicted HRV" unit="ms" />
        </div>

        <div className="chart-panel">
          <PanelHeader title="Sleep quality" meta="score" />
          <LineChart points={sleepSeries} actualLabel="Actual sleep" predictedLabel="Predicted sleep" unit="" />
        </div>

        <div className="chart-panel">
          <PanelHeader title="Fatigue" meta="score" />
          <LineChart points={fatigueSeries} actualLabel="Actual fatigue" predictedLabel="Predicted fatigue" unit="" />
        </div>
      </section>

      <section className="split-section">
        <div className="model-panel">
          <PanelHeader
            title="Model trace"
            meta={model ? `${model.model_id}@${model.model_version}` : "no model"}
          />
          {model ? (
            <>
              <dl className="trace-list">
                <div>
                  <dt>Input window</dt>
                  <dd>
                    {model.input_start_date} to {model.input_end_date}
                  </dd>
                </div>
                <div>
                  <dt>Feature count</dt>
                  <dd>{model.feature_set.length}</dd>
                </div>
                <div>
                  <dt>Generated</dt>
                  <dd>{formatDateTime(data.generated_at)}</dd>
                </div>
              </dl>
              <TagList items={model.feature_set} />
              <TextList title="Assumptions" items={model.assumptions} />
            </>
          ) : (
            <p className="muted">No model output is available.</p>
          )}
        </div>

        <div className="privacy-panel">
          <PanelHeader title="Public data boundary" meta="redacted export" />
          <div className="privacy-grid">
            <BoundaryItem label="Shown" value="Daily aggregates, model metadata, prediction errors" />
            <BoundaryItem label="Omitted" value="Notes, weight, location, raw timelines, credentials" />
            <BoundaryItem label="Current source" value={latestRecord?.source_id ?? "synthetic sample"} />
            <BoundaryItem label="Next gate" value="Validated device adapters before real equipment" />
          </div>
        </div>
      </section>

      <section className="table-section">
        <PanelHeader title="Recent evaluations" meta="prediction vs actual" />
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Date</th>
                <th>Recovery</th>
                <th>HRV error</th>
                <th>RHR error</th>
                <th>Fatigue error</th>
                <th>Sleep error</th>
                <th>Likely drivers</th>
              </tr>
            </thead>
            <tbody>
              {data.evaluations.slice(-8).map((item) => (
                <tr key={item.target_date}>
                  <td>{item.target_date}</td>
                  <td>
                    <StatePill state={item.prediction.recovery_state} />
                  </td>
                  <td>{formatSigned(item.hrv_error)}</td>
                  <td>{formatSigned(item.rhr_error)}</td>
                  <td>{formatSigned(item.fatigue_error)}</td>
                  <td>{formatSigned(item.sleep_error)}</td>
                  <td>{item.prediction.hints.join("; ")}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section className="roadmap-band">
        <PanelHeader title="Open roadmap" meta="equipment-ready foundation" />
        <div className="roadmap-grid">
          <RoadmapItem title="Public dashboard" detail="React Vite view over redacted JSON exports." />
          <RoadmapItem title="Adapter registry" detail="Apple Health, Oura, Garmin, Fitbit, and CSV contracts." />
          <RoadmapItem title="Model registry" detail="Versioned predictions, assumptions, and comparable error reports." />
          <RoadmapItem title="Research simulations" detail="Biology topics stay as safe toy models until real approvals exist." />
        </div>
      </section>
    </main>
  );
}

function StatusScreen({ title, detail }: { title: string; detail: string }) {
  return (
    <main className="status-screen">
      <div>
        <img className="status-logo" src={logoUrl} alt="Flyto2 logo" />
        <p className="eyebrow">Flyto2</p>
        <h1>{title}</h1>
        <p>{detail}</p>
      </div>
    </main>
  );
}

function MetricCard({
  label,
  value,
  unit,
  tone,
}: {
  label: string;
  value: string;
  unit?: string;
  tone?: string;
}) {
  return (
    <article className={`metric-card ${tone ? `tone-${tone}` : ""}`}>
      <span>{label}</span>
      <strong>
        {value}
        {unit ? <small>{unit}</small> : null}
      </strong>
    </article>
  );
}

function PanelHeader({ title, meta }: { title: string; meta: string }) {
  return (
    <div className="panel-header">
      <h2>{title}</h2>
      <span>{meta}</span>
    </div>
  );
}

function LineChart({
  points,
  actualLabel,
  predictedLabel,
  unit,
}: {
  points: SeriesPoint[];
  actualLabel: string;
  predictedLabel: string;
  unit: string;
}) {
  const width = 720;
  const height = 260;
  const padding = 34;
  const allValues = points.flatMap((point) => [point.actual, point.predicted]);
  const min = Math.min(...allValues);
  const max = Math.max(...allValues);
  const span = Math.max(1, max - min);

  const toX = (index: number) =>
    padding + (index * (width - padding * 2)) / Math.max(1, points.length - 1);
  const toY = (value: number) => height - padding - ((value - min) * (height - padding * 2)) / span;
  const pathFor = (key: "actual" | "predicted") =>
    points.map((point, index) => `${index === 0 ? "M" : "L"} ${toX(index)} ${toY(point[key])}`).join(" ");

  return (
    <div className="chart-box">
      <svg viewBox={`0 0 ${width} ${height}`} role="img" aria-label={`${actualLabel} and ${predictedLabel}`}>
        <line className="axis" x1={padding} y1={height - padding} x2={width - padding} y2={height - padding} />
        <line className="axis" x1={padding} y1={padding} x2={padding} y2={height - padding} />
        <text className="axis-label" x={padding} y={padding - 10}>
          {formatNumber(max)}
          {unit ? ` ${unit}` : ""}
        </text>
        <text className="axis-label" x={padding} y={height - 8}>
          {formatNumber(min)}
          {unit ? ` ${unit}` : ""}
        </text>
        <path className="line actual" d={pathFor("actual")} />
        <path className="line predicted" d={pathFor("predicted")} />
        {points.map((point, index) => (
          <g key={`${point.label}-${index}`}>
            <circle className="dot actual-dot" cx={toX(index)} cy={toY(point.actual)} r="3.5" />
            <circle className="dot predicted-dot" cx={toX(index)} cy={toY(point.predicted)} r="3.5" />
          </g>
        ))}
      </svg>
      <div className="chart-legend">
        <span className="legend-actual">{actualLabel}</span>
        <span className="legend-predicted">{predictedLabel}</span>
      </div>
    </div>
  );
}

function TagList({ items }: { items: string[] }) {
  return (
    <div className="tag-list">
      {items.map((item) => (
        <span key={item}>{item.replaceAll("_", " ")}</span>
      ))}
    </div>
  );
}

function TextList({ title, items }: { title: string; items: string[] }) {
  return (
    <div className="text-list">
      <h3>{title}</h3>
      <ul>
        {items.map((item) => (
          <li key={item}>{item}</li>
        ))}
      </ul>
    </div>
  );
}

function BoundaryItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="boundary-item">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function RoadmapItem({ title, detail }: { title: string; detail: string }) {
  return (
    <article className="roadmap-item">
      <h3>{title}</h3>
      <p>{detail}</p>
    </article>
  );
}

function StatePill({ state }: { state: string }) {
  return <span className={`state-pill tone-${state}`}>{formatState(state)}</span>;
}

function formatState(value?: string) {
  if (!value) {
    return "Unknown";
  }
  return value
    .split("_")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function formatNumber(value: number) {
  return value.toFixed(1);
}

function formatSigned(value: number) {
  return `${value >= 0 ? "+" : ""}${value.toFixed(1)}`;
}

function formatDateTime(value: string) {
  return new Intl.DateTimeFormat("en", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(new Date(value));
}
