.PHONY: format

format:
	@echo "Formatting markdown files..."
	# MD013 (Line length): 技術文書ではコマンドやURLで長い行が発生するため無効化
	# MD033 (No inline HTML): Mermaid図やテーブル内の改行に <br> を使用するため無効化
	# MD010 (No hard tabs): Goコードブロック内で idiomatic な Tab インデントを維持するため無効化
	npx markdownlint "**/*.md" --fix --disable MD013 MD033 MD010
