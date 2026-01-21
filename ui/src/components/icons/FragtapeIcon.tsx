import { ComponentProps } from "react";

type Props = {
  animated?: boolean;
} & ComponentProps<"svg">

const bigStrokeWidth = 4
const smallStrokeWidth = 2

export const FragtapeIcon = ({ animated = false, ...props }: Props) => {
  if (animated) return <FragtapeIconAnimated {...props} />

  return <FragtapeIconNormal {...props} />
}

const FragtapeIconNormal = (props: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="-5 -5 169 169"
      {...props}
    >
      <g
        stroke="currentColor"
        strokeWidth={bigStrokeWidth}
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      >
        <circle cx="72" cy="72" r="40" />

        <line x1="72" y1="22" x2="72" y2="52" />
        <line x1="72" y1="92" x2="72" y2="122" />

        <line x1="22" y1="72" x2="52" y2="72" />
        <line x1="92" y1="72" x2="122" y2="72" />
      </g>

      <g
        stroke="currentColor"
        strokeWidth={smallStrokeWidth}
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      >
        <line x1="72" y1="62" x2="72" y2="82" />
        <line x1="62" y1="72" x2="82" y2="72" />
      </g>

      <g
        stroke="currentColor"
        strokeWidth={bigStrokeWidth}
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      >
        <path d="
        M 141 82
        a 70 70 0 0 0 0 -20
        l -15 -1
        a 55 55 0 0 0 -9 -20
        l 10 -11
        a 70 70 0 0 0 -14 -14
        l -10 11
        a 55 55 0 0 0 -20 -9
        l -1 -15
        a 70 70 0 0 0 -20 0
        l -1 15
        a 55 55 0 0 0 -21 9
        l -11 -11
        a 70 70 0 0 0 -14 14
        l 11 11
        a 55 55 0 0 0 -8 20
        l -15 1
        a 70 70 0 0 0 0 20
        l 15 1
        a 55 55 0 0 0 9 20
        l -11 11
        a 70 70 0 0 0 14 14
        l 11 -11
        a 55 55 0 0 0 20 9
        l 1 15
        a 70 70 0 0 0 20 0
        l 1 -15
        a 55 55 0 0 0 21 -9
        l 10 11
        a 70 70 0 0 0 14 -14
        l -11 -11
        a 55 55 0 0 0 9 -21
        Z
      " />
      </g>
    </svg>
  );
}

const FragtapeIconAnimated = (props: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="-5 -5 169 169"
      {...props}
    >
      <g
        stroke="currentColor"
        strokeWidth={bigStrokeWidth}
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      >
        <circle cx="72" cy="72" r="40" />

        <line x1="72" y1="22" x2="72" y2="52" />
        <line x1="72" y1="92" x2="72" y2="122" />

        <line x1="22" y1="72" x2="52" y2="72" />
        <line x1="92" y1="72" x2="122" y2="72" />
      </g>

      <g
        stroke="currentColor"
        strokeWidth={smallStrokeWidth}
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      >
        <line x1="72" y1="62" x2="72" y2="82" />
        <line x1="62" y1="72" x2="82" y2="72" />
      </g>

      <g
        stroke="currentColor"
        strokeWidth={bigStrokeWidth}
        strokeLinecap="round"
        strokeLinejoin="round"
        fill="none"
      >
        <g>
          <animateTransform
            attributeName="transform"
            type="translate"
            additive="sum"
            dur="1.6s"
            repeatCount="indefinite"
            values="0 0; -6.48 -6.48; -6.48 -6.48; 0 0"
            keyTimes="0; 0.4; 0.6; 1"
            calcMode="spline"
            keySplines="0.2 0 0.2 1; 0 0 1 1; 0.2 0 0.2 1"
          />

          <animateTransform
            attributeName="transform"
            type="scale"
            additive="sum"
            dur="1.6s"
            repeatCount="indefinite"
            values="1; 1.09; 1.09; 1"
            keyTimes="0; 0.4; 0.6; 1"
            calcMode="spline"
            keySplines="0.2 0 0.2 1; 0 0 1 1; 0.2 0 0.2 1"
          />

          <animateTransform
            attributeName="transform"
            type="rotate"
            additive="sum"
            dur="1.6s"
            repeatCount="indefinite"
            values="0 72 72; 0 72 72; 90 72 72; 90 72 72"
            keyTimes="0; 0.32; 0.78; 1"
            calcMode="spline"
            keySplines="0 0 1 1; 0.2 0 0.2 1; 0 0 1 1"
          />

          <path d="
        M 141 82
        a 70 70 0 0 0 0 -20
        l -15 -1
        a 55 55 0 0 0 -9 -20
        l 10 -11
        a 70 70 0 0 0 -14 -14
        l -10 11
        a 55 55 0 0 0 -20 -9
        l -1 -15
        a 70 70 0 0 0 -20 0
        l -1 15
        a 55 55 0 0 0 -21 9
        l -11 -11
        a 70 70 0 0 0 -14 14
        l 11 11
        a 55 55 0 0 0 -8 20
        l -15 1
        a 70 70 0 0 0 0 20
        l 15 1
        a 55 55 0 0 0 9 20
        l -11 11
        a 70 70 0 0 0 14 14
        l 11 -11
        a 55 55 0 0 0 20 9
        l 1 15
        a 70 70 0 0 0 20 0
        l 1 -15
        a 55 55 0 0 0 21 -9
        l 10 11
        a 70 70 0 0 0 14 -14
        l -11 -11
        a 55 55 0 0 0 9 -21
        Z
      " />
        </g>
      </g>
    </svg>
  )
}
