import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@solidjs/testing-library';
import userEvent from '@testing-library/user-event';
import { BoxTrack } from './BoxTrack';

function getFilledBoxCount(): number {
  const buttons = screen.getAllByRole('button');
  return buttons.filter((btn) => btn.className.includes('boxFilled')).length;
}

describe('BoxTrack', () => {
  describe('Box Filling', () => {
    it('renders correct number of boxes', () => {
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={0}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      expect(buttons).toHaveLength(9);
    });

    it('fills boxes up to the value', () => {
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={5}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      expect(getFilledBoxCount()).toBe(5);
    });

    it('shows no filled boxes when value is 0', () => {
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={0}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      expect(getFilledBoxCount()).toBe(0);
    });

    it('fills all boxes when value equals totalBoxes', () => {
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={9}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      expect(getFilledBoxCount()).toBe(9);
    });
  });

  describe('Click Behavior', () => {
    it('clicking empty box fills up to that box', async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={0}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      await user.click(buttons[4]); // Click 5th box (index 4)

      expect(onChange).toHaveBeenCalledWith(5);
    });

    it('clicking the currently filled box decrements value', async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={5}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      await user.click(buttons[4]); // Click 5th box which is the last filled

      expect(onChange).toHaveBeenCalledWith(4);
    });

    it('clicking a box higher than current value fills up to that box', async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={3}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      await user.click(buttons[7]); // Click 8th box (index 7)

      expect(onChange).toHaveBeenCalledWith(8);
    });

    it('clicking a box lower than current value fills up to that box', async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={7}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      await user.click(buttons[2]); // Click 3rd box (index 2)

      expect(onChange).toHaveBeenCalledWith(3);
    });

    it('clicking first box when value is 1 sets value to 0', async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={1}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      await user.click(buttons[0]); // Click 1st box

      expect(onChange).toHaveBeenCalledWith(0);
    });

    it('clicking last box fills all boxes', async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={0}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      await user.click(buttons[8]); // Click 9th box (index 8)

      expect(onChange).toHaveBeenCalledWith(9);
    });
  });

  describe('Accented Boxes', () => {
    it('applies accent style to specified boxes', () => {
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={0}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[3, 6, 9]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      const accentedButtons = buttons.filter((btn) =>
        btn.className.includes('boxAccented')
      );

      expect(accentedButtons).toHaveLength(3);
      expect(buttons[2].className).toContain('boxAccented'); // 3rd box
      expect(buttons[5].className).toContain('boxAccented'); // 6th box
      expect(buttons[8].className).toContain('boxAccented'); // 9th box
    });

    it('non-accented boxes do not have accent style', () => {
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={0}
          onChange={onChange}
          label="Test"
          totalBoxes={9}
          accentedBoxes={[3, 9]}
        />
      ));

      const buttons = screen.getAllByRole('button');
      expect(buttons[0].className).not.toContain('boxAccented');
      expect(buttons[1].className).not.toContain('boxAccented');
      expect(buttons[3].className).not.toContain('boxAccented'); // 4th box
    });
  });

  describe('Label', () => {
    it('displays the label text', () => {
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={0}
          onChange={onChange}
          label="Courage"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      expect(screen.getByText('Courage')).toBeInTheDocument();
    });

    it('shows current value in tooltip', () => {
      const onChange = vi.fn();
      render(() => (
        <BoxTrack
          value={5}
          onChange={onChange}
          label="Courage"
          totalBoxes={9}
          accentedBoxes={[]}
        />
      ));

      const labelWrapper = screen.getByText('Courage').parentElement;
      expect(labelWrapper).toHaveAttribute('data-tip', '5');
    });
  });
});
