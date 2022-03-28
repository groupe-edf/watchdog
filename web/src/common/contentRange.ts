export interface Range {
  unit: string;
  start?: number | null;
  end?: number | null;
  size?: number | null;
}

const parseContentRange = (input: string): Range | null => {
  const matches = input.match(/^(\w+) ((\d+)-(\d+)|\*)\/(\d+|\*)$/);
  if (!matches) return null;
  const [, unit, , start, end, size] = matches;
  const range = {
    unit,
    start: start != null ? Number(start) : null,
    end: end != null ? Number(end) : null,
    size: size === "*" ? null : Number(size),
  };
  if (range.start === null && range.end === null && range.size === null)
    return null;
  return range;
}

export default parseContentRange
